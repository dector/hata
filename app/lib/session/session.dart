class Session {
  const Session({
    required this.token,
    required this.serverUrl,
    required this.username,
    required this.validUntil,
    required this.createdAt,
    this.displayName,
  });

  final String token;
  final String serverUrl;
  final String username;
  final String? displayName;
  final DateTime validUntil;
  final DateTime createdAt;

  bool get isExpired => !validUntil.isAfter(DateTime.now());
}
