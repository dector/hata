import 'package:flutter/material.dart';

import '../theme/hata_colors.dart';

enum ConnectionStatus { unknown, online, offline }

class HataTopBar extends StatelessWidget {
  const HataTopBar({
    super.key,
    required this.homeName,
    required this.connectionStatus,
    required this.isSyncing,
    required this.displayName,
    required this.username,
    this.hasUnreadNotifications = false,
    this.onNotifications,
    this.onProfile,
    this.onHomeSelector,
  });

  final String homeName;
  final ConnectionStatus connectionStatus;
  final bool isSyncing;
  final String? displayName;
  final String username;
  final bool hasUnreadNotifications;
  final VoidCallback? onNotifications;
  final VoidCallback? onProfile;
  final VoidCallback? onHomeSelector;

  @override
  Widget build(BuildContext context) {
    final statusLabel = isSyncing ? 'Syncing…' : _statusLabel(connectionStatus);
    final avatarLetter = _avatarLetter(displayName ?? username);

    return Padding(
      padding: const EdgeInsets.all(16),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Flexible(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                InkWell(
                  borderRadius: BorderRadius.circular(12),
                  onTap: onHomeSelector,
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Flexible(
                        child: Text(
                          homeName,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                          style: const TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      const Icon(Icons.keyboard_arrow_down),
                    ],
                  ),
                ),
                const SizedBox(height: 4),
                Row(
                  children: [
                    _StatusDot(connectionStatus: connectionStatus),
                    const SizedBox(width: 6),
                    Text(
                      statusLabel,
                      style: const TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          Row(
            children: [
              IconButton(
                tooltip: 'Notifications',
                onPressed: onNotifications,
                icon: hasUnreadNotifications
                    ? Badge(
                        smallSize: 8,
                        backgroundColor: HataColors.errorVariant,
                        child: const Icon(Icons.notifications),
                      )
                    : const Icon(Icons.notifications),
              ),
              const SizedBox(width: 4),
              Semantics(
                label: 'Profile',
                button: true,
                child: InkWell(
                  borderRadius: BorderRadius.circular(24),
                  onTap: onProfile,
                  child: CircleAvatar(
                    radius: 24,
                    backgroundColor: HataColors.avatarBackground,
                    child: Text(
                      avatarLetter,
                      style: const TextStyle(
                        color: HataColors.avatarIcon,
                        fontSize: 20,
                        fontWeight: FontWeight.w600,
                      ),
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

  static String _statusLabel(ConnectionStatus status) {
    switch (status) {
      case ConnectionStatus.online:
        return 'Online';
      case ConnectionStatus.offline:
        return 'Offline';
      case ConnectionStatus.unknown:
        return 'Unknown';
    }
  }

  static String _avatarLetter(String value) {
    final trimmed = value.trim();
    if (trimmed.isEmpty) return '?';
    return trimmed.characters.first.toUpperCase();
  }
}

class _StatusDot extends StatelessWidget {
  const _StatusDot({required this.connectionStatus});

  final ConnectionStatus connectionStatus;

  @override
  Widget build(BuildContext context) {
    final color = switch (connectionStatus) {
      ConnectionStatus.online => HataColors.primary,
      ConnectionStatus.offline => HataColors.errorVariant,
      ConnectionStatus.unknown => HataColors.onSurfaceVariant,
    };

    return Container(
      width: 8,
      height: 8,
      decoration: BoxDecoration(color: color, shape: BoxShape.circle),
    );
  }
}
