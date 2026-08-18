import 'package:flutter/material.dart';

import '../api/api_error.dart';
import '../api/hata_api_client.dart';
import '../models/api_device.dart';
import '../models/home_device.dart';
import '../models/house.dart';
import '../session/session.dart';
import '../theme/hata_colors.dart';
import '../widgets/device_card.dart';
import '../widgets/hata_top_bar.dart';
import '../widgets/home_fab_menu.dart';
import 'shopping_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({
    super.key,
    required this.session,
    required this.apiClient,
    required this.onSessionInvalid,
  });

  final Session session;
  final HataApiClient apiClient;
  final Future<void> Function() onSessionInvalid;

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> with WidgetsBindingObserver {
  bool _showFabMenu = false;
  List<House> _houses = const [];
  List<HomeDevice> _devices = const [];
  ConnectionStatus _connectionStatus = ConnectionStatus.unknown;
  bool _isInitialLoading = true;
  bool _isRefreshing = false;
  String? _errorMessage;
  final Set<String> _pendingDeviceIds = <String>{};

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    _loadHome(initial: true);
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state == AppLifecycleState.resumed && !_isInitialLoading) {
      _loadHome();
    }
  }

  Future<void> _loadHome({bool initial = false}) async {
    if (initial) {
      setState(() {
        _isInitialLoading = true;
        _errorMessage = null;
      });
    } else {
      setState(() {
        _isRefreshing = true;
        _errorMessage = null;
      });
    }

    try {
      final houses = await widget.apiClient.fetchHouses(
        widget.session.serverUrl,
        widget.session.token,
      );
      final apiDevices = await widget.apiClient.fetchDevices(
        widget.session.serverUrl,
        widget.session.token,
      );
      if (!mounted) return;
      setState(() {
        _houses = houses;
        _devices = apiDevices
            .map((device) => _toHomeDevice(device, houses: houses))
            .toList();
        _connectionStatus = ConnectionStatus.online;
        _isInitialLoading = false;
        _isRefreshing = false;
        _errorMessage = null;
      });
    } catch (error) {
      if (_isUnauthorized(error)) {
        await widget.onSessionInvalid();
        return;
      }
      if (!mounted) return;
      setState(() {
        _connectionStatus = ConnectionStatus.offline;
        _isInitialLoading = false;
        _isRefreshing = false;
        _errorMessage = _messageFor(error);
      });
    }
  }

  HomeDevice _toHomeDevice(ApiDevice device, {required List<House> houses}) {
    final state = device.state.toLowerCase();
    final availability = device.availability.toLowerCase();
    final isOnline = availability == 'online';
    final isKnownPowerState = state == 'on' || state == 'off';
    final displayAsOn = isOnline && state == 'on';
    final isLight =
        device.capabilities.light ||
        device.integration.id.toLowerCase().contains('wiz') ||
        device.integration.id.toLowerCase().contains('light');

    return HomeDevice(
      id: device.id,
      houseId: device.houseId,
      title: device.name.isEmpty ? device.id : device.name,
      status: _statusText(device, houses: houses),
      icon: isLight ? Icons.lightbulb : Icons.devices_other,
      isOn: displayAsOn,
      canControl: isOnline && isKnownPowerState,
      rawState: device.state,
      availability: device.availability,
      isAwaitingConfirmation: _pendingDeviceIds.contains(device.id),
    );
  }

  String _statusText(ApiDevice device, {required List<House> houses}) {
    final parts = <String>[];
    final integrationId = device.integration.id;
    if (integrationId.isNotEmpty && integrationId.toLowerCase() != 'unknown') {
      parts.add(integrationId);
    }

    final availability = device.availability.toLowerCase();
    if (availability == 'offline') {
      parts.add('Offline');
    } else {
      final state = device.state.toLowerCase();
      parts.add(switch (state) {
        'on' => 'On',
        'off' => 'Off',
        _ => 'Unknown',
      });
    }

    if (houses.length > 1) {
      final house = houses
          .where((house) => house.id == device.houseId)
          .firstOrNull;
      parts.add(house?.displayName ?? device.houseId);
    }

    return parts.join(' • ');
  }

  Future<void> _toggleDevice(String id, bool value) async {
    final index = _devices.indexWhere((device) => device.id == id);
    if (index == -1) return;

    final device = _devices[index];
    if (!device.canControl || _pendingDeviceIds.contains(id)) return;

    final previousDevices = List<HomeDevice>.from(_devices);
    _pendingDeviceIds.add(id);
    setState(() {
      _devices = List<HomeDevice>.from(_devices)
        ..[index] = device.copyWith(
          isOn: value,
          rawState: value ? 'on' : 'off',
          status: _optimisticStatus(device, value),
          isAwaitingConfirmation: true,
        );
    });

    try {
      final result = await widget.apiClient.setDeviceState(
        widget.session.serverUrl,
        widget.session.token,
        device.houseId,
        device.id,
        value,
      );
      if (!mounted) return;
      final updatedIndex = _devices.indexWhere((item) => item.id == id);
      _pendingDeviceIds.remove(id);
      if (updatedIndex == -1) {
        setState(() {});
        return;
      }
      final updated = _devices[updatedIndex];
      final state = result.state.toLowerCase();
      setState(() {
        _connectionStatus = ConnectionStatus.online;
        _devices = List<HomeDevice>.from(_devices)
          ..[updatedIndex] = updated.copyWith(
            isOn: state == 'on',
            rawState: result.state,
            status: _optimisticStatus(updated, state == 'on'),
            isAwaitingConfirmation: false,
          );
      });
      unawaitedRefresh();
    } catch (error) {
      if (_isUnauthorized(error)) {
        await widget.onSessionInvalid();
        return;
      }
      if (!mounted) return;
      _pendingDeviceIds.remove(id);
      setState(() {
        _devices = previousDevices;
        _connectionStatus = ConnectionStatus.offline;
        _errorMessage = _messageFor(error);
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(_errorMessage ?? 'Device update failed')),
      );
    }
  }

  void unawaitedRefresh() {
    Future<void>.delayed(const Duration(milliseconds: 250), () {
      if (mounted) _loadHome();
    });
  }

  String _optimisticStatus(HomeDevice device, bool isOn) {
    final parts = device.status.split(' • ');
    final stateIndex = parts.indexWhere(
      (part) => part == 'On' || part == 'Off' || part == 'Unknown',
    );
    if (stateIndex == -1) {
      parts.add(isOn ? 'On' : 'Off');
    } else {
      parts[stateIndex] = isOn ? 'On' : 'Off';
    }
    return parts.join(' • ');
  }

  bool _isUnauthorized(Object error) =>
      error is ApiError && (error.statusCode == 401 || error.statusCode == 403);

  String _messageFor(Object error) =>
      error is ApiError ? error.message : 'Could not reach the server';

  void _openShopping() {
    setState(() => _showFabMenu = false);
    Navigator.of(
      context,
    ).push(MaterialPageRoute<void>(builder: (_) => const ShoppingScreen()));
  }

  @override
  Widget build(BuildContext context) {
    final homeName = _houses.isEmpty ? 'My Home' : _houses.first.displayName;
    final isSyncing = _isRefreshing || _pendingDeviceIds.isNotEmpty;

    return Scaffold(
      floatingActionButton: HomeFabMenu(
        isExpanded: _showFabMenu,
        onToggle: () => setState(() => _showFabMenu = !_showFabMenu),
        onShopping: _openShopping,
      ),
      body: SafeArea(
        child: Column(
          children: [
            HataTopBar(
              homeName: homeName,
              connectionStatus: _connectionStatus,
              isSyncing: isSyncing,
              displayName: widget.session.displayName,
              username: widget.session.username,
            ),
            Expanded(
              child: RefreshIndicator(
                color: HataColors.primary,
                backgroundColor: HataColors.surface,
                onRefresh: () => _loadHome(),
                child: CustomScrollView(
                  physics: const AlwaysScrollableScrollPhysics(),
                  slivers: [
                    _WelcomeHeader(
                      displayName:
                          widget.session.displayName ?? widget.session.username,
                    ),
                    ..._contentSlivers(),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  List<Widget> _contentSlivers() {
    if (_isInitialLoading) {
      return const [
        SliverFillRemaining(
          hasScrollBody: false,
          child: Center(child: CircularProgressIndicator()),
        ),
      ];
    }

    if (_devices.isEmpty && _errorMessage != null) {
      return [
        SliverFillRemaining(
          hasScrollBody: false,
          child: _MessageState(
            icon: Icons.cloud_off,
            title: 'Home is offline',
            message: _errorMessage!,
            actionLabel: 'Retry',
            onAction: () => _loadHome(initial: true),
          ),
        ),
      ];
    }

    if (_devices.isEmpty) {
      return const [
        SliverFillRemaining(
          hasScrollBody: false,
          child: _MessageState(
            icon: Icons.devices_other,
            title: 'No devices yet',
            message: 'Devices visible to your account will appear here.',
          ),
        ),
      ];
    }

    return [
      if (_errorMessage != null)
        SliverToBoxAdapter(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
            child: Text(
              _errorMessage!,
              style: const TextStyle(color: HataColors.errorVariant),
            ),
          ),
        ),
      SliverPadding(
        padding: const EdgeInsets.fromLTRB(16, 8, 16, 96),
        sliver: SliverLayoutBuilder(
          builder: (context, constraints) {
            final width = constraints.crossAxisExtent;
            final columns = width >= 920
                ? 4
                : width >= 640
                ? 3
                : 2;

            return SliverGrid(
              delegate: SliverChildBuilderDelegate((context, index) {
                final device = _devices[index];
                return DeviceCard(
                  device: device,
                  onChanged: (value) => _toggleDevice(device.id, value),
                );
              }, childCount: _devices.length),
              gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: columns,
                mainAxisSpacing: 16,
                crossAxisSpacing: 16,
                mainAxisExtent: 192,
              ),
            );
          },
        ),
      ),
    ];
  }
}

class _WelcomeHeader extends StatelessWidget {
  const _WelcomeHeader({required this.displayName});

  final String displayName;

  @override
  Widget build(BuildContext context) {
    final firstName = displayName.trim().split(RegExp(r'\s+')).first;
    return SliverToBoxAdapter(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(16, 16, 16, 24),
        child: Text(
          'Welcome home,\n$firstName',
          style: const TextStyle(
            fontSize: 36,
            height: 1.16,
            fontWeight: FontWeight.w800,
          ),
        ),
      ),
    );
  }
}

class _MessageState extends StatelessWidget {
  const _MessageState({
    required this.icon,
    required this.title,
    required this.message,
    this.actionLabel,
    this.onAction,
  });

  final IconData icon;
  final String title;
  final String message;
  final String? actionLabel;
  final VoidCallback? onAction;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, size: 48, color: HataColors.onSurfaceVariant),
          const SizedBox(height: 16),
          Text(
            title,
            textAlign: TextAlign.center,
            style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w700),
          ),
          const SizedBox(height: 8),
          Text(
            message,
            textAlign: TextAlign.center,
            style: const TextStyle(color: HataColors.onSurfaceVariant),
          ),
          if (actionLabel != null && onAction != null) ...[
            const SizedBox(height: 20),
            FilledButton(onPressed: onAction, child: Text(actionLabel!)),
          ],
        ],
      ),
    );
  }
}
