class AuthSession {
  const AuthSession({
    required this.token,
    required this.validUntil,
    required this.displayName,
  });

  final String token;
  final DateTime validUntil;
  final String? displayName;

  factory AuthSession.fromJson(Map<String, dynamic> json) {
    final session = json['session'] as Map<String, dynamic>? ?? const {};
    final user = json['user'] as Map<String, dynamic>? ?? const {};
    final validUntilValue = session['validUntil'] as String?;

    return AuthSession(
      token: session['token'] as String? ?? '',
      validUntil: validUntilValue == null
          ? DateTime.fromMillisecondsSinceEpoch(0)
          : DateTime.parse(validUntilValue),
      displayName: user['displayName'] as String?,
    );
  }
}
