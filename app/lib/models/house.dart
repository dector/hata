class House {
  const House({required this.id, required this.displayName, required this.role});

  final String id;
  final String displayName;
  final String role;

  factory House.fromJson(Map<String, dynamic> json) {
    final id = json['id']?.toString() ?? '';
    return House(
      id: id,
      displayName: json['displayName']?.toString() ?? id,
      role: json['role']?.toString() ?? '',
    );
  }
}
