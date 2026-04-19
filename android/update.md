# Dependency Update Plan

## 1) Build toolchain
- `agp`: `9.0.0` → `9.1.1`
- `kotlin`: `2.3.10` → `2.3.20`
- `ksp`: `2.3.4` → `2.3.6`
- `gradle/plugins/build-config/build.gradle.kts`:
  - `com.android.tools.build:gradle:9.0.0` → `9.1.1`

## 2) Compose / UI stack
- `compose` BOM: `2026.01.01` → `2026.03.01`
- `androidx-activity-compose`: `1.12.3` → `1.13.0`
- `androidx-core`: `1.17.0` → `1.18.0`

## 3) Navigation
- `androidx-navigation3`: `1.0.0` → `1.1.0`

## 4) Data / persistence
- `androidx-room`: `2.8.0` → `2.8.4`

## 5) Serialization / DI ✅
- `kotlinx-serialization`: `1.10.0` → `1.11.0` _(updated)_
- `hilt`: `2.59.1` → `2.59.2` _(updated)_
