import 'package:flutter/material.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key, required this.onLoginSuccess});

  final VoidCallback onLoginSuccess;

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _serverController = TextEditingController(text: 'http://10.0.2.2:8080');
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _isServerConnected = false;

  @override
  void dispose() {
    _serverController.dispose();
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  bool get _canConnect => _serverController.text.trim().isNotEmpty;
  bool get _canLogin =>
      _usernameController.text.trim().isNotEmpty &&
      _passwordController.text.isNotEmpty;

  void _connect() {
    if (!_canConnect) return;
    setState(() => _isServerConnected = true);
  }

  void _changeServer() {
    setState(() => _isServerConnected = false);
  }

  void _login() {
    if (!_canLogin) return;
    widget.onLoginSuccess();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 420),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  const Text(
                    'Hata',
                    textAlign: TextAlign.center,
                    style: TextStyle(fontSize: 42, fontWeight: FontWeight.w800),
                  ),
                  const SizedBox(height: 32),
                  if (_isServerConnected)
                    _CredentialsForm(
                      usernameController: _usernameController,
                      passwordController: _passwordController,
                      onChanged: () => setState(() {}),
                      onLogin: _login,
                      onChangeServer: _changeServer,
                      canLogin: _canLogin,
                    )
                  else
                    _ServerForm(
                      controller: _serverController,
                      onChanged: () => setState(() {}),
                      onConnect: _connect,
                      canConnect: _canConnect,
                    ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _ServerForm extends StatelessWidget {
  const _ServerForm({
    required this.controller,
    required this.onChanged,
    required this.onConnect,
    required this.canConnect,
  });

  final TextEditingController controller;
  final VoidCallback onChanged;
  final VoidCallback onConnect;
  final bool canConnect;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        TextField(
          controller: controller,
          onChanged: (_) => onChanged(),
          keyboardType: TextInputType.url,
          decoration: const InputDecoration(
            labelText: 'Server URL',
            border: OutlineInputBorder(),
          ),
        ),
        const SizedBox(height: 16),
        FilledButton(
          onPressed: canConnect ? onConnect : null,
          child: const Text('Connect'),
        ),
      ],
    );
  }
}

class _CredentialsForm extends StatelessWidget {
  const _CredentialsForm({
    required this.usernameController,
    required this.passwordController,
    required this.onChanged,
    required this.onLogin,
    required this.onChangeServer,
    required this.canLogin,
  });

  final TextEditingController usernameController;
  final TextEditingController passwordController;
  final VoidCallback onChanged;
  final VoidCallback onLogin;
  final VoidCallback onChangeServer;
  final bool canLogin;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const Text(
          'Welcome to Hata',
          textAlign: TextAlign.center,
          style: TextStyle(fontSize: 26, fontWeight: FontWeight.w700),
        ),
        const SizedBox(height: 24),
        TextField(
          controller: usernameController,
          onChanged: (_) => onChanged(),
          decoration: const InputDecoration(
            labelText: 'Username',
            border: OutlineInputBorder(),
          ),
        ),
        const SizedBox(height: 12),
        TextField(
          controller: passwordController,
          onChanged: (_) => onChanged(),
          obscureText: true,
          decoration: const InputDecoration(
            labelText: 'Password',
            border: OutlineInputBorder(),
          ),
        ),
        const SizedBox(height: 16),
        FilledButton(
          onPressed: canLogin ? onLogin : null,
          child: const Text('Login'),
        ),
        TextButton(
          onPressed: onChangeServer,
          child: const Text('Change Server'),
        ),
      ],
    );
  }
}
