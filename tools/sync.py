#!/usr/bin/env -S uv run --script
# /// script
# requires-python = ">=3.11"
# ///
"""One-way sync from a remote git worktree into this local repo.

Usage:
  ./tools/sync.py pi@factory:~/hata
  ./tools/sync.py --debounce 2.5 --log-file .sync-from-factory.log pi@factory:~/hata

The script must be started from the local repository root.  It watches the
remote tree with inotifywait, then rsyncs only files reported by remote
`git ls-files --cached --others --exclude-standard`.
"""

from __future__ import annotations

import argparse
import hashlib
import os
import select
import shlex
import subprocess
import sys
import time
from datetime import datetime
from pathlib import Path


class SyncError(RuntimeError):
    pass


# inotifywait watches the tree before rsync's git-based filtering happens.
# Exclude noisy generated directories up front so Android/Flutter builds don't
# constantly retrigger syncs.  Keep this as a regex because that is what
# inotifywait --exclude accepts.
WATCH_EXCLUDE_REGEX = (
    r"(^|/)("
    r"\.git|\.jj|"
    r"build|\.dart_tool|\.gradle|\.cxx|"
    r"out|target|node_modules|"
    r"\.pub-cache|\.symlinks|ephemeral|Pods|DerivedData"
    r")(/|$)"
)


def log(message: str, log_file: Path | None = None) -> None:
    line = f"[{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] {message}"
    print(line, flush=True)
    if log_file is not None:
        log_file.parent.mkdir(parents=True, exist_ok=True)
        with log_file.open("a", encoding="utf-8") as f:
            f.write(line + "\n")


def parse_source(source: str) -> tuple[str, str]:
    if ":" not in source:
        raise SyncError("source must look like 'user@host:/path' or 'user@host:~/path'")
    host, remote_dir = source.split(":", 1)
    if not host or not remote_dir:
        raise SyncError("source must include both host and path, e.g. 'pi@factory:~/hata'")
    return host, remote_dir


def shell_cd_path(path: str) -> str:
    """Return a shell expression suitable after `cd` on the remote host.

    Keep common ~/... paths unquoted so the remote shell expands them.
    Other paths are shell-quoted.
    """
    if path == "~" or (path.startswith("~/") and "'" not in path and not any(c.isspace() for c in path)):
        return path
    return shlex.quote(path)


def run(cmd: list[str], *, input_bytes: bytes | None = None) -> subprocess.CompletedProcess[bytes]:
    return subprocess.run(cmd, input=input_bytes, stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=False)


def remote_pwd(host: str, remote_dir: str) -> str:
    remote_cmd = f"cd {shell_cd_path(remote_dir)} && pwd -P"
    cp = run(["ssh", host, "bash", "-lc", shlex.quote(remote_cmd)])
    if cp.returncode != 0:
        raise SyncError(cp.stderr.decode(errors="replace").strip() or "failed to resolve remote path")
    resolved = cp.stdout.decode(errors="replace").strip()
    if not resolved.startswith("/"):
        raise SyncError(f"remote path did not resolve to an absolute path: {resolved}")
    return resolved


def remote_git_manifest(host: str, remote_dir: str) -> bytes:
    remote_cmd = (
        f"cd {shell_cd_path(remote_dir)} && "
        "git ls-files -z --cached --others --exclude-standard"
    )
    # ssh concatenates remote argv with spaces before passing it to the
    # remote user's shell. Quote the bash -lc payload so it stays one argument;
    # otherwise only `cd` becomes the -c command and the rest runs elsewhere.
    cp = run(["ssh", host, "bash", "-lc", shlex.quote(remote_cmd)])
    if cp.returncode != 0:
        raise SyncError(cp.stderr.decode(errors="replace").strip() or "failed to read remote git manifest")
    return cp.stdout


def manifest_entries(manifest: bytes) -> set[str]:
    return {item.decode("utf-8", "surrogateescape") for item in manifest.split(b"\0") if item}


def remove_deleted_files(old_manifest: bytes, new_manifest: bytes, local_dir: Path, log_file: Path | None) -> None:
    old = manifest_entries(old_manifest)
    new = manifest_entries(new_manifest)

    for rel in sorted(old - new):
        rel_path = Path(rel)
        if rel_path.is_absolute() or ".." in rel_path.parts:
            log(f"skip suspicious deleted path: {rel}", log_file)
            continue

        path = local_dir / rel_path
        try:
            if path.is_file() or path.is_symlink():
                path.unlink()
                log(f"deleted {rel}", log_file)

                parent = path.parent
                while parent != local_dir and local_dir in parent.parents:
                    try:
                        parent.rmdir()
                    except OSError:
                        break
                    parent = parent.parent
        except OSError as exc:
            log(f"failed to delete {rel}: {exc}", log_file)


def rsync_from_manifest(source: str, manifest_file: Path, local_dir: Path) -> None:
    remote_source = source.rstrip("/") + "/"
    cp = run([
        "rsync",
        "-az",
        "--from0",
        f"--files-from={manifest_file}",
        remote_source,
        str(local_dir) + "/",
    ])
    if cp.returncode != 0:
        raise SyncError(cp.stderr.decode(errors="replace").strip() or "rsync failed")


def sync_once(source: str, host: str, remote_dir: str, state_dir: Path, local_dir: Path, log_file: Path | None) -> None:
    new_manifest = remote_git_manifest(host, remote_dir)
    old_manifest_path = state_dir / "manifest.prev"
    new_manifest_path = state_dir / "manifest.new"

    new_manifest_path.write_bytes(new_manifest)

    log("sync started", log_file)
    rsync_from_manifest(source, new_manifest_path, local_dir)

    if old_manifest_path.exists():
        remove_deleted_files(old_manifest_path.read_bytes(), new_manifest, local_dir, log_file)

    new_manifest_path.replace(old_manifest_path)
    log("sync finished", log_file)


def start_remote_watcher(host: str, remote_dir: str) -> subprocess.Popen[str]:
    remote_cmd = f"""
cd {shell_cd_path(remote_dir)} || exit
command -v inotifywait >/dev/null || {{
  echo 'inotifywait missing; install inotify-tools on the remote host' >&2
  exit 127
}}
inotifywait -m -r \
  -e close_write,create,delete,move,attrib \
  --exclude {shlex.quote(WATCH_EXCLUDE_REGEX)} \
  --format '%w%f' .
""".strip()
    return subprocess.Popen(
        ["ssh", host, "bash", "-lc", shlex.quote(remote_cmd)],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=1,
    )


def drain_stderr(proc: subprocess.Popen[str]) -> str:
    if proc.stderr is None:
        return ""
    chunks: list[str] = []
    while True:
        ready, _, _ = select.select([proc.stderr], [], [], 0)
        if not ready:
            break
        line = proc.stderr.readline()
        if not line:
            break
        chunks.append(line.rstrip())
    return "\n".join(chunks)


def watch_loop(source: str, host: str, remote_dir: str, state_dir: Path, local_dir: Path, debounce: float, log_file: Path | None) -> None:
    proc = start_remote_watcher(host, remote_dir)
    assert proc.stdout is not None

    log(f"watching {source} with {debounce:g}s debounce", log_file)
    log(f"watch exclude regex: {WATCH_EXCLUDE_REGEX}", log_file)

    try:
        while True:
            if proc.poll() is not None:
                err = drain_stderr(proc)
                raise SyncError(f"remote watcher exited with code {proc.returncode}" + (f":\n{err}" if err else ""))

            ready, _, _ = select.select([proc.stdout], [], [], 0.5)
            if not ready:
                err = drain_stderr(proc)
                if err:
                    log(err, log_file)
                continue

            first = proc.stdout.readline()
            if not first:
                continue
            changed_count = 1
            last_change = first.strip()
            deadline = time.monotonic() + debounce

            while True:
                timeout = max(0.0, deadline - time.monotonic())
                ready, _, _ = select.select([proc.stdout], [], [], timeout)
                if not ready:
                    break
                line = proc.stdout.readline()
                if not line:
                    break
                changed_count += 1
                last_change = line.strip()
                deadline = time.monotonic() + debounce

            log(f"remote changes settled: {changed_count} event(s), last={last_change}", log_file)
            sync_once(source, host, remote_dir, state_dir, local_dir, log_file)
    finally:
        proc.terminate()


def main() -> int:
    parser = argparse.ArgumentParser(description="Watch a remote repo and sync non-gitignored files back here.")
    parser.add_argument("source", help="remote source, e.g. pi@factory:~/hata")
    parser.add_argument("--debounce", type=float, default=1.0, help="seconds to wait after the last event before syncing (default: 1.0)")
    parser.add_argument("--log-file", type=Path, default=None, help="optional file to append timestamped logs to")
    args = parser.parse_args()

    if args.debounce < 0:
        parser.error("--debounce must be >= 0")

    local_dir = Path.cwd()
    if not (local_dir / ".git").exists():
        parser.error("run this script from the repository root")

    try:
        host, remote_dir_input = parse_source(args.source)
        remote_dir = remote_pwd(host, remote_dir_input)
        resolved_source = f"{host}:{remote_dir}"
        source_id = hashlib.sha256(resolved_source.encode()).hexdigest()[:16]
        state_dir = local_dir / ".git" / "sync-from-remote" / source_id
        state_dir.mkdir(parents=True, exist_ok=True)

        log(f"resolved source: {resolved_source}", args.log_file)
        log("initial sync", args.log_file)
        sync_once(resolved_source, host, remote_dir, state_dir, local_dir, args.log_file)
        watch_loop(resolved_source, host, remote_dir, state_dir, local_dir, args.debounce, args.log_file)
    except KeyboardInterrupt:
        log("stopped", args.log_file)
        return 130
    except SyncError as exc:
        log(f"error: {exc}", args.log_file)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
