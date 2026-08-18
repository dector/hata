import 'package:flutter/material.dart';

import '../theme/hata_colors.dart';

class HomeFabMenu extends StatelessWidget {
  const HomeFabMenu({
    super.key,
    required this.isExpanded,
    required this.onToggle,
    required this.onShopping,
  });

  final bool isExpanded;
  final VoidCallback onToggle;
  final VoidCallback onShopping;

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        if (isExpanded) ...[
          Material(
            color: HataColors.surface,
            borderRadius: BorderRadius.circular(20),
            clipBehavior: Clip.antiAlias,
            child: InkWell(
              onTap: onShopping,
              child: const SizedBox(
                width: 220,
                child: Padding(
                  padding: EdgeInsets.symmetric(horizontal: 20, vertical: 16),
                  child: Row(
                    children: [
                      _ShoppingMenuDot(),
                      SizedBox(width: 14),
                      Text(
                        'Shopping',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
          const SizedBox(height: 16),
        ],
        FloatingActionButton(
          onPressed: onToggle,
          child: AnimatedRotation(
            turns: isExpanded ? 0.5 : 0,
            duration: const Duration(milliseconds: 150),
            child: const Icon(Icons.keyboard_arrow_up),
          ),
        ),
      ],
    );
  }
}

class _ShoppingMenuDot extends StatelessWidget {
  const _ShoppingMenuDot();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 10,
      height: 10,
      decoration: const BoxDecoration(
        color: Color(0xFF4CAF50),
        shape: BoxShape.circle,
      ),
    );
  }
}
