class ShoppingItem {
  const ShoppingItem(this.id, this.name, {this.isChecked = false});

  final String id;
  final String name;
  final bool isChecked;

  ShoppingItem copyWith({bool? isChecked}) {
    return ShoppingItem(id, name, isChecked: isChecked ?? this.isChecked);
  }
}
