import 'package:flutter/material.dart';

import 'api/api.dart';
import 'screens/home_screen.dart';
import 'screens/login_screen.dart';
import 'session/session.dart';
import 'session/session_database.dart';
import 'session/session_repository.dart';

class AppShell extends StatefulWidget {
  AppShell({
    super.key,
    Api? api,
    SessionRepository? sessionRepository,
  }) : api = api ?? Api(),
       sessionRepository = sessionRepository ?? SqliteSessionRepository();

  final Api api;
  final SessionRepository sessionRepository;

  @override
  State<AppShell> createState() => _AppShellState();
}

class _AppShellState extends State<AppShell> {
  Session? _session;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadSession();
  }

  Future<void> _loadSession() async {
    final session = await widget.sessionRepository.load();
    if (!mounted) return;

    if (session == null || session.isExpired) {
      await widget.sessionRepository.clear();
      if (mounted) {
        setState(() => _isLoading = false);
      }
      return;
    }

    try {
      await widget.api.fetchHouse(session.serverUrl, session.token);
      if (mounted) {
        setState(() {
          _session = session;
          _isLoading = false;
        });
      }
    } catch (_) {
      await widget.sessionRepository.clear();
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  void _startSession(Session session) {
    setState(() {
      _session = session;
      _isLoading = false;
    });
  }

  Future<void> _clearSession() async {
    await widget.sessionRepository.clear();
    if (!mounted) return;
    setState(() {
      _session = null;
      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    return _session != null
        ? HomeScreen(
            session: _session!,
            apiClient: widget.api,
            onSessionInvalid: _clearSession,
          )
        : LoginScreen(
            apiClient: widget.api,
            sessionRepository: widget.sessionRepository,
            onLoginSuccess: _startSession,
          );
  }
}
