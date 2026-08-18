import 'package:flutter/material.dart';

import '../models/home_device.dart';
import '../theme/hata_colors.dart';

class DeviceCard extends StatelessWidget {
  const DeviceCard({super.key, required this.device, required this.onChanged});

  final HomeDevice device;
  final ValueChanged<bool> onChanged;

  @override
  Widget build(BuildContext context) {
    final isOffline = device.availability.toLowerCase() == 'offline';
    final isActive = device.isOn && !isOffline;
    final cardColor = isActive ? HataColors.primary : HataColors.surface;
    final foregroundColor = Colors.white;
    final secondaryTextColor = isActive
        ? Colors.white.withValues(alpha: 0.85)
        : HataColors.onSurfaceVariant;
    final iconColor = isActive ? Colors.white : HataColors.onSurfaceVariant;
    final iconBackground = isActive
        ? HataColors.primaryContainer
        : HataColors.surfaceVariant;

    return Semantics(
      label: '${device.title}, ${device.rawState}',
      child: Opacity(
        opacity: isOffline ? 0.55 : 1,
        child: Card(
          color: cardColor,
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(24),
          ),
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
                    if (device.isAwaitingConfirmation)
                      SizedBox(
                        width: 48,
                        height: 48,
                        child: Padding(
                          padding: const EdgeInsets.all(12),
                          child: CircularProgressIndicator(
                            strokeWidth: 2.5,
                            color: isActive ? Colors.white : HataColors.primary,
                          ),
                        ),
                      )
                    else
                      Switch(
                        value: device.isOn,
                        activeThumbColor: Colors.white,
                        activeTrackColor: HataColors.primaryContainerVariant,
                        inactiveThumbColor: HataColors.onSurfaceVariant,
                        inactiveTrackColor: HataColors.surfaceVariant,
                        onChanged: device.canControl ? onChanged : null,
                      ),
                  ],
                ),
                const Spacer(),
                Text(
                  device.title,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(
                    color: foregroundColor,
                    fontSize: 20,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  device.status,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(color: secondaryTextColor, fontSize: 14),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
