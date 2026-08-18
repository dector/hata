import 'package:flutter/material.dart';

import '../models/home_device.dart';
import '../theme/hata_colors.dart';

class DeviceCard extends StatelessWidget {
  const DeviceCard({super.key, required this.device, required this.onChanged});

  final HomeDevice device;
  final ValueChanged<bool> onChanged;

  @override
  Widget build(BuildContext context) {
    final cardColor = device.isOn ? HataColors.primary : HataColors.surface;
    final iconColor = device.isOn ? Colors.white : HataColors.onSurfaceVariant;
    final iconBackground = device.isOn
        ? HataColors.primaryContainer
        : HataColors.surfaceVariant;

    return Semantics(
      label: '${device.title}, ${device.isOn ? 'on' : 'off'}',
      child: Card(
        color: cardColor,
        elevation: 0,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
        clipBehavior: Clip.antiAlias,
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  CircleAvatar(
                    radius: 24,
                    backgroundColor: iconBackground,
                    child: Icon(device.icon, color: iconColor),
                  ),
                  Switch(
                    value: device.isOn,
                    activeThumbColor: Colors.white,
                    activeTrackColor: HataColors.primaryContainerVariant,
                    inactiveThumbColor: HataColors.onSurfaceVariant,
                    inactiveTrackColor: HataColors.surfaceVariant,
                    onChanged: onChanged,
                  ),
                ],
              ),
              const Spacer(),
              Text(
                device.title,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w500,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                device.isOn
                    ? 'On • ${device.status}'
                    : 'Off • ${device.status}',
                style: TextStyle(
                  color: Colors.white.withValues(alpha: 0.85),
                  fontSize: 14,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
