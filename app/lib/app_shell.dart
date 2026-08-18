import 'package:flutter/material.dart';

import 'screens/home_screen.dart';
import 'screens/login_screen.dart';

class AppShell extends StatefulWidget {
  const AppShell({super.key});

  @override
  State<AppShell> createState() => _AppShellState();
}

class _AppShellState extends State<AppShell> {
  bool _hasSession = false;

  void _startSession() {
    setState(() => _hasSession = true);
  }

  @override
  Widget build(BuildContext context) {
    return _hasSession
        ? const HomeScreen()
        : LoginScreen(onLoginSuccess: _startSession);
  }
}
