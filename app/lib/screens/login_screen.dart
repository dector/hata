import 'package:flutter/material.dart';

import '../api/api_error.dart';
import '../api/hata_api_client.dart';
import '../models/server_info.dart';
import '../session/session.dart';
import '../session/session_repository.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({
    super.key,
    required this.apiClient,
    required this.sessionRepository,
    required this.onLoginSuccess,
  });

  final HataApiClient apiClient;
  final SessionRepository sessionRepository;
  final ValueChanged<Session> onLoginSuccess;

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _serverController = TextEditingController(text: 'http://10.0.2.2:4501');
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();
  ServerInfo? _serverInfo;
  bool _isConnecting = false;
  bool _isLoggingIn = false;
  String? _error;

  @override
  void dispose() {
    _serverController.dispose();
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  bool get _isServerConnected => _serverInfo != null;
  bool get _canConnect =>
      !_isConnecting && _serverController.text.trim().isNotEmpty;
  bool get _canLogin =>
      !_isLoggingIn &&
      _usernameController.text.trim().isNotEmpty &&
      _passwordController.text.isNotEmpty;

  Future<void> _connect() async {
    if (!_canConnect) return;
    setState(() {
      _isConnecting = true;
      _error = null;
    });

    try {
      final serverInfo = await widget.apiClient.ping(_serverController.text);
      if (mounted) {
        setState(() => _serverInfo = serverInfo);
      }
    } catch (error) {
      if (mounted) {
        setState(
          () => _error = _messageFor(error, 'Could not connect to server'),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isConnecting = false);
      }
    }
  }

  void _changeServer() {
    setState(() {
      _serverInfo = null;
      _error = null;
    });
  }

  Future<void> _login() async {
    if (!_canLogin) return;
    setState(() {
      _isLoggingIn = true;
      _error = null;
    });

    try {
      final serverUrl = widget.apiClient.normalizeBaseUrl(
        _serverController.text,
      );
      final username = _usernameController.text.trim();
      final auth = await widget.apiClient.login(
        serverUrl,
        username,
        _passwordController.text,
      );
      await widget.apiClient.fetchHouse(serverUrl, auth.token);
      final session = Session(
        token: auth.token,
        serverUrl: serverUrl,
        username: username,
        displayName: auth.displayName,
        validUntil: auth.validUntil,
        createdAt: DateTime.now(),
      );
      await widget.sessionRepository.save(session);
      if (mounted) {
        widget.onLoginSuccess(session);
      }
    } catch (error) {
      if (mounted) {
        setState(() => _error = _messageFor(error, 'Login failed'));
      }
    } finally {
      if (mounted) {
        setState(() => _isLoggingIn = false);
      }
    }
  }

  String _messageFor(Object error, String fallback) {
    if (error is ApiError) return error.message;
    return fallback;
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
                  if (_error != null) ...[
                    Text(
                      _error!,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        color: Theme.of(context).colorScheme.error,
                      ),
                    ),
                    const SizedBox(height: 16),
                  ],
                  if (_isServerConnected)
                    _CredentialsForm(
                      serverName: _serverInfo!.serverName,
                      usernameController: _usernameController,
                      passwordController: _passwordController,
                      onChanged: () => setState(() {}),
                      onLogin: _login,
                      onChangeServer: _changeServer,
                      canLogin: _canLogin,
                      isLoggingIn: _isLoggingIn,
                    )
                  else
                    _ServerForm(
                      controller: _serverController,
                      onChanged: () => setState(() {}),
                      onConnect: _connect,
                      canConnect: _canConnect,
                      isConnecting: _isConnecting,
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
    required this.isConnecting,
  });

  final TextEditingController controller;
  final VoidCallback onChanged;
  final VoidCallback onConnect;
  final bool canConnect;
  final bool isConnecting;

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
          child: Text(isConnecting ? 'Connecting…' : 'Connect'),
        ),
      ],
    );
  }
}

class _CredentialsForm extends StatelessWidget {
  const _CredentialsForm({
    required this.serverName,
    required this.usernameController,
    required this.passwordController,
    required this.onChanged,
    required this.onLogin,
    required this.onChangeServer,
    required this.canLogin,
    required this.isLoggingIn,
  });

  final String serverName;
  final TextEditingController usernameController;
  final TextEditingController passwordController;
  final VoidCallback onChanged;
  final VoidCallback onLogin;
  final VoidCallback onChangeServer;
  final bool canLogin;
  final bool isLoggingIn;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(
          serverName,
          textAlign: TextAlign.center,
          style: const TextStyle(fontSize: 26, fontWeight: FontWeight.w700),
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
          child: Text(isLoggingIn ? 'Logging in…' : 'Login'),
        ),
        TextButton(
          onPressed: isLoggingIn ? null : onChangeServer,
          child: const Text('Change Server'),
        ),
      ],
    );
  }
}
