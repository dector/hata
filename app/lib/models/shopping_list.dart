class ShoppingList {
  const ShoppingList({
    required this.uid,
    required this.name,
    required this.createdAt,
    required this.updatedAt,
  });

  final String uid;
  final String name;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  String get id => uid;

  factory ShoppingList.fromJson(Map<String, dynamic> json) {
    return ShoppingList(
      uid: json['uid']?.toString() ?? '',
      name: json['name']?.toString() ?? '',
      createdAt: _parseDate(json['createdAt']),
      updatedAt: _parseDate(json['updatedAt']),
    );
  }
}

DateTime? _parseDate(Object? value) {
  if (value == null) return null;
  return DateTime.tryParse(value.toString());
}
