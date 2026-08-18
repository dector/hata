import 'package:flutter/material.dart';

import '../models/shopping_item.dart';
import '../theme/hata_colors.dart';

class ShoppingTile extends StatelessWidget {
  const ShoppingTile({
    super.key,
    required this.item,
    required this.onChanged,
    this.enabled = true,
  });

  final ShoppingItem item;
  final ValueChanged<bool?> onChanged;
  final bool enabled;

  @override
  Widget build(BuildContext context) {
    return Card(
      color: HataColors.surface,
      elevation: 0,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
      child: CheckboxListTile(
        value: item.isChecked,
        onChanged: enabled ? onChanged : null,
        activeColor: HataColors.primary,
        checkboxShape: const CircleBorder(),
        title: Text(
          item.name,
          style: TextStyle(
            fontSize: 16,
            decoration: item.isChecked ? TextDecoration.lineThrough : null,
            color: item.isChecked
                ? Colors.white.withValues(alpha: 0.55)
                : Colors.white,
          ),
        ),
      ),
    );
  }
}
