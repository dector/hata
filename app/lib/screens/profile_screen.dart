import 'package:flutter/material.dart';
import 'package:package_info_plus/package_info_plus.dart';

import '../api/api.dart';
import '../models/server_info.dart';
import '../session/session.dart';
import '../theme/hata_colors.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({
    super.key,
    required this.session,
    required this.apiClient,
    required this.onLogout,
  });

  final Session session;
  final Api apiClient;
  final Future<void> Function() onLogout;

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  PackageInfo? _packageInfo;
  ServerInfo? _serverInfo;

  @override
  void initState() {
    super.initState();
    _loadPackageInfo();
    _loadServerInfo();
  }

  Future<void> _loadPackageInfo() async {
    try {
      final packageInfo = await PackageInfo.fromPlatform();
      if (!mounted) return;
      setState(() => _packageInfo = packageInfo);
    } catch (_) {
      // Keep fallback values visible if platform package info is unavailable.
    }
  }

  Future<void> _loadServerInfo() async {
    try {
      final serverInfo = await widget.apiClient.ping(widget.session.serverUrl);
      if (!mounted) return;
      setState(() => _serverInfo = serverInfo);
    } catch (_) {
      // Keep fallback values visible if the server cannot be reached.
    }
  }

  Future<void> _showLogoutDialog() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => const _LogoutDialog(),
    );
    if (confirmed != true || !mounted) return;

    await widget.onLogout();
    if (!mounted) return;
    Navigator.of(context).popUntil((route) => route.isFirst);
  }

  @override
  Widget build(BuildContext context) {
    final session = widget.session;
    final userName = (session.displayName?.trim().isNotEmpty ?? false)
        ? session.displayName!.trim()
        : session.username;
    final packageInfo = _packageInfo;

    return Scaffold(
      backgroundColor: HataColors.background,
      appBar: AppBar(
        backgroundColor: HataColors.background,
        foregroundColor: Colors.white,
        title: const Text('Profile'),
      ),
      body: SafeArea(
        child: Stack(
          children: [
            SingleChildScrollView(
              padding: const EdgeInsets.fromLTRB(32, 16, 32, 120),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  const _UserAvatar(),
                  const SizedBox(height: 24),
                  Text(
                    userName,
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 32,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                  const SizedBox(height: 48),
                  _ProfileInfoSections(
                    serverUrl: session.serverUrl,
                    serverVersion: _serverInfo?.version ?? '—',
                    accountStatus: session.isExpired ? 'Inactive' : 'Active',
                    appPackage: packageInfo?.packageName ?? '—',
                    appVersion: packageInfo?.version ?? '—',
                    sessionCreatedAt: session.createdAt,
                  ),
                ],
              ),
            ),
            Align(
              alignment: Alignment.bottomCenter,
              child: Padding(
                padding: const EdgeInsets.fromLTRB(32, 0, 32, 24),
                child: _LogoutButton(onPressed: _showLogoutDialog),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _UserAvatar extends StatelessWidget {
  const _UserAvatar();

  @override
  Widget build(BuildContext context) {
    return const CircleAvatar(
      radius: 60,
      backgroundColor: HataColors.avatarBackground,
      child: Icon(Icons.person, color: HataColors.avatarIcon, size: 64),
    );
  }
}

class _ProfileInfoSections extends StatelessWidget {
  const _ProfileInfoSections({
    required this.serverUrl,
    required this.serverVersion,
    required this.accountStatus,
    required this.appVersion,
    required this.appPackage,
    required this.sessionCreatedAt,
  });

  final String serverUrl;
  final String serverVersion;
  final String accountStatus;
  final String appVersion;
  final String appPackage;
  final DateTime sessionCreatedAt;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const _SectionTitle('User'),
        _InfoCard(
          child: _InfoRow(
            label: 'Created',
            value: _formatSessionCreatedAt(sessionCreatedAt),
          ),
        ),
        const SizedBox(height: 20),
        const _SectionTitle('Server'),
        _InfoCard(
          child: Column(
            children: [
              _InfoRow(label: 'URL', value: serverUrl),
              const SizedBox(height: 12),
              _InfoRow(label: 'Version', value: serverVersion),
              const SizedBox(height: 12),
              _StatusRow(label: 'Status', status: accountStatus),
            ],
          ),
        ),
        const SizedBox(height: 20),
        const _SectionTitle('App'),
        _InfoCard(
          child: Column(
            children: [
              _InfoRow(label: 'Package', value: appPackage),
              const SizedBox(height: 12),
              _InfoRow(label: 'Version', value: appVersion),
            ],
          ),
        ),
      ],
    );
  }
}

class _SectionTitle extends StatelessWidget {
  const _SectionTitle(this.title);

  final String title;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Text(
        title,
        style: const TextStyle(
          color: HataColors.onSurfaceVariant,
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}

class _InfoCard extends StatelessWidget {
  const _InfoCard({required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: EdgeInsets.zero,
      color: HataColors.surface,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(padding: const EdgeInsets.all(24), child: child),
    );
  }
}

class _InfoRow extends StatelessWidget {
  const _InfoRow({required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            color: HataColors.onSurfaceVariant,
            fontSize: 16,
          ),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: Text(
            value,
            textAlign: TextAlign.end,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.w500,
            ),
          ),
        ),
      ],
    );
  }
}

class _StatusRow extends StatelessWidget {
  const _StatusRow({required this.label, required this.status});

  final String label;
  final String status;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Text(
          label,
          style: const TextStyle(
            color: HataColors.onSurfaceVariant,
            fontSize: 16,
          ),
        ),
        const Spacer(),
        Container(
          width: 8,
          height: 8,
          decoration: BoxDecoration(
            color: _statusColor(status),
            shape: BoxShape.circle,
          ),
        ),
        const SizedBox(width: 8),
        Text(
          status,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 16,
            fontWeight: FontWeight.w500,
          ),
        ),
      ],
    );
  }
}

class _LogoutButton extends StatelessWidget {
  const _LogoutButton({required this.onPressed});

  final VoidCallback onPressed;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      height: 56,
      child: FilledButton(
        onPressed: onPressed,
        style: FilledButton.styleFrom(
          backgroundColor: HataColors.error,
          foregroundColor: HataColors.onPrimary,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
        child: const Text(
          'Logout',
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600),
        ),
      ),
    );
  }
}

class _LogoutDialog extends StatelessWidget {
  const _LogoutDialog();

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      backgroundColor: HataColors.surface,
      title: const Text('Logout', style: TextStyle(color: Colors.white)),
      content: const Text(
        'Are you sure you want to logout?',
        style: TextStyle(color: HataColors.onSurfaceVariant),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(true),
          child: const Text(
            'Logout',
            style: TextStyle(color: HataColors.error),
          ),
        ),
        TextButton(
          onPressed: () => Navigator.of(context).pop(false),
          child: const Text('Cancel', style: TextStyle(color: Colors.white)),
        ),
      ],
    );
  }
}

String _formatSessionCreatedAt(DateTime value) {
  final local = value.toLocal();
  final month = const [
    'Jan',
    'Feb',
    'Mar',
    'Apr',
    'May',
    'Jun',
    'Jul',
    'Aug',
    'Sep',
    'Oct',
    'Nov',
    'Dec',
  ][local.month - 1];
  final hour = local.hour.toString().padLeft(2, '0');
  final minute = local.minute.toString().padLeft(2, '0');
  return '$month ${local.day}, ${local.year} • $hour:$minute';
}

Color _statusColor(String status) {
  switch (status.toLowerCase()) {
    case 'active':
      return HataColors.primary;
    case 'inactive':
      return HataColors.error;
    case 'warning':
    case 'pending':
    default:
      return const Color(0xFFF2C94C);
  }
}
