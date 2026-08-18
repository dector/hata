class ServerInfo {
  const ServerInfo({
    required this.serverName,
    required this.version,
    required this.apiVersion,
  });

  final String serverName;
  final String version;
  final String apiVersion;

  factory ServerInfo.fromJson(Map<String, dynamic> json) {
    return ServerInfo(
      serverName: json['serverName'] as String? ?? 'Hata Server',
      version: json['version'] as String? ?? '',
      apiVersion: json['apiVersion'] as String? ?? 'latest',
    );
  }
}
