class ShoppingItem {
  const ShoppingItem({
    required this.uid,
    required this.name,
    required this.position,
    required this.checkedAt,
    required this.checkedByUserId,
    required this.deletedAt,
    required this.createdAt,
    required this.updatedAt,
  });

  final String uid;
  final String name;
  final int position;
  final DateTime? checkedAt;
  final int? checkedByUserId;
  final DateTime? deletedAt;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  String get id => uid;
  bool get isChecked => checkedAt != null;

  factory ShoppingItem.fromJson(Map<String, dynamic> json) {
    return ShoppingItem(
      uid: json['uid']?.toString() ?? '',
      name: json['name']?.toString() ?? '',
      position: _parseInt(json['position']),
      checkedAt: _parseDate(json['checkedAt']),
      checkedByUserId: _parseNullableInt(json['checkedByUserId']),
      deletedAt: _parseDate(json['deletedAt']),
      createdAt: _parseDate(json['createdAt']),
      updatedAt: _parseDate(json['updatedAt']),
    );
  }
}

int _parseInt(Object? value) => _parseNullableInt(value) ?? 0;

int? _parseNullableInt(Object? value) {
  if (value == null) return null;
  if (value is int) return value;
  return int.tryParse(value.toString());
}

DateTime? _parseDate(Object? value) {
  if (value == null) return null;
  return DateTime.tryParse(value.toString());
}
