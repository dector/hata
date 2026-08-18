import 'package:flutter/material.dart';

import '../models/shopping_item.dart';
import '../theme/hata_colors.dart';
import '../widgets/shopping_tile.dart';

class ShoppingScreen extends StatefulWidget {
  const ShoppingScreen({super.key});

  @override
  State<ShoppingScreen> createState() => _ShoppingScreenState();
}

class _ShoppingScreenState extends State<ShoppingScreen> {
  final List<ShoppingItem> _items = [
    ShoppingItem('1', 'Coffee beans'),
    ShoppingItem('2', 'Milk'),
    ShoppingItem('3', 'Fresh bread'),
    ShoppingItem('4', 'Tomatoes', isChecked: true),
  ];

  void _toggleItem(String id, bool? checked) {
    setState(() {
      final index = _items.indexWhere((item) => item.id == id);
      _items[index] = _items[index].copyWith(isChecked: checked ?? false);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: HataColors.background,
        foregroundColor: Colors.white,
        title: const Text('Shopping'),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {},
        child: const Icon(Icons.add),
      ),
      body: SafeArea(
        child: RefreshIndicator(
          color: HataColors.primary,
          backgroundColor: HataColors.surface,
          onRefresh: () async {
            await Future<void>.delayed(const Duration(milliseconds: 500));
          },
          child: ListView.separated(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.fromLTRB(16, 16, 16, 96),
            itemBuilder: (context, index) {
              final item = _items[index];
              return ShoppingTile(
                item: item,
                onChanged: (value) => _toggleItem(item.id, value),
              );
            },
            separatorBuilder: (context, index) => const SizedBox(height: 12),
            itemCount: _items.length,
          ),
        ),
      ),
    );
  }
}
