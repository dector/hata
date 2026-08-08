import 'package:flutter/material.dart';

void main() {
  runApp(const HataApp());
}

class HataApp extends StatelessWidget {
  const HataApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Hata',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: HataColors.primary,
          brightness: Brightness.dark,
        ),
        scaffoldBackgroundColor: HataColors.background,
        useMaterial3: true,
      ),
      home: const HomeScreen(),
    );
  }
}

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  bool _showFabMenu = false;
  final List<HomeDevice> _devices = [
    HomeDevice(
      id: 'cooler',
      title: 'Air Cooler',
      status: 'Living room',
      icon: Icons.ac_unit,
      isOn: true,
    ),
    HomeDevice(
      id: 'kitchen',
      title: 'Kitchen Light',
      status: 'Kitchen',
      icon: Icons.lightbulb,
      isOn: false,
    ),
    HomeDevice(
      id: 'bedroom',
      title: 'Bedroom Light',
      status: 'Bedroom',
      icon: Icons.lightbulb_outline,
      isOn: true,
    ),
    HomeDevice(
      id: 'garage',
      title: 'Garage Door',
      status: 'Garage',
      icon: Icons.keyboard_arrow_up,
      isOn: false,
    ),
  ];

  void _toggleDevice(String id, bool value) {
    setState(() {
      final index = _devices.indexWhere((device) => device.id == id);
      _devices[index] = _devices[index].copyWith(isOn: value);
    });
  }

  void _openShopping() {
    setState(() => _showFabMenu = false);
    Navigator.of(
      context,
    ).push(MaterialPageRoute<void>(builder: (_) => const ShoppingScreen()));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      floatingActionButton: _HomeFabMenu(
        isExpanded: _showFabMenu,
        onToggle: () => setState(() => _showFabMenu = !_showFabMenu),
        onShopping: _openShopping,
      ),
      body: SafeArea(
        child: Column(
          children: [
            const HataTopBar(),
            Expanded(
              child: RefreshIndicator(
                color: HataColors.primary,
                backgroundColor: HataColors.surface,
                onRefresh: () async {
                  await Future<void>.delayed(const Duration(milliseconds: 500));
                },
                child: CustomScrollView(
                  physics: const AlwaysScrollableScrollPhysics(),
                  slivers: [
                    const SliverToBoxAdapter(child: _WelcomeHeader()),
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
                            delegate: SliverChildBuilderDelegate((
                              context,
                              index,
                            ) {
                              final device = _devices[index];
                              return DeviceCard(
                                device: device,
                                onChanged: (value) =>
                                    _toggleDevice(device.id, value),
                              );
                            }, childCount: _devices.length),
                            gridDelegate:
                                SliverGridDelegateWithFixedCrossAxisCount(
                                  crossAxisCount: columns,
                                  mainAxisSpacing: 16,
                                  crossAxisSpacing: 16,
                                  mainAxisExtent: 192,
                                ),
                          );
                        },
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class HataColors {
  static const primary = Color(0xFF8FB899);
  static const onPrimary = Color(0xFF1E2423);
  static const background = Color(0xFF1E2423);
  static const surface = Color(0xFF2D3533);
  static const surfaceVariant = Color(0xFF3D4845);
  static const onSurfaceVariant = Color(0xFF6B7875);
  static const primaryContainer = Color(0xFFA8C5B0);
  static const primaryContainerVariant = Color(0xFFC8DFD0);
  static const errorVariant = Color(0xFFEF5350);
  static const avatarBackground = Color(0xFFB39B8D);
  static const avatarIcon = Color(0xFFE8D4C4);
}

class HomeDevice {
  const HomeDevice({
    required this.id,
    required this.title,
    required this.status,
    required this.icon,
    required this.isOn,
  });

  final String id;
  final String title;
  final String status;
  final IconData icon;
  final bool isOn;

  HomeDevice copyWith({bool? isOn}) {
    return HomeDevice(
      id: id,
      title: title,
      status: status,
      icon: icon,
      isOn: isOn ?? this.isOn,
    );
  }
}

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

class _WelcomeHeader extends StatelessWidget {
  const _WelcomeHeader();

  @override
  Widget build(BuildContext context) {
    return const Padding(
      padding: EdgeInsets.fromLTRB(16, 16, 16, 24),
      child: Text(
        'Welcome home,\nDan',
        style: TextStyle(
          fontSize: 36,
          height: 1.16,
          fontWeight: FontWeight.w800,
        ),
      ),
    );
  }
}

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

class _HomeFabMenu extends StatelessWidget {
  const _HomeFabMenu({
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

class ShoppingItem {
  const ShoppingItem(this.id, this.name, {this.isChecked = false});

  final String id;
  final String name;
  final bool isChecked;

  ShoppingItem copyWith({bool? isChecked}) {
    return ShoppingItem(id, name, isChecked: isChecked ?? this.isChecked);
  }
}

class ShoppingTile extends StatelessWidget {
  const ShoppingTile({super.key, required this.item, required this.onChanged});

  final ShoppingItem item;
  final ValueChanged<bool?> onChanged;

  @override
  Widget build(BuildContext context) {
    return Card(
      color: HataColors.surface,
      elevation: 0,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
      child: CheckboxListTile(
        value: item.isChecked,
        onChanged: onChanged,
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
