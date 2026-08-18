import 'package:flutter/material.dart';

import '../theme/hata_colors.dart';

class HataTopBar extends StatelessWidget {
  const HataTopBar({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          const Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Text(
                    'My Home',
                    style: TextStyle(fontSize: 24, fontWeight: FontWeight.w700),
                  ),
                  SizedBox(width: 8),
                  Icon(Icons.keyboard_arrow_down),
                ],
              ),
              SizedBox(height: 4),
              Row(
                children: [
                  _StatusDot(),
                  SizedBox(width: 6),
                  Text(
                    'Online',
                    style: TextStyle(fontSize: 12, fontWeight: FontWeight.w500),
                  ),
                ],
              ),
            ],
          ),
          Row(
            children: [
              IconButton(
                tooltip: 'Notifications',
                onPressed: () {},
                icon: Badge(
                  smallSize: 8,
                  backgroundColor: HataColors.errorVariant,
                  child: const Icon(Icons.notifications),
                ),
              ),
              const SizedBox(width: 4),
              Semantics(
                label: 'Profile',
                button: true,
                child: const CircleAvatar(
                  radius: 24,
                  backgroundColor: HataColors.avatarBackground,
                  child: Text(
                    'D',
                    style: TextStyle(
                      color: HataColors.avatarIcon,
                      fontSize: 20,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _StatusDot extends StatelessWidget {
  const _StatusDot();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 8,
      height: 8,
      decoration: const BoxDecoration(
        color: HataColors.primary,
        shape: BoxShape.circle,
      ),
    );
  }
}
