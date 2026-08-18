import 'package:flutter/material.dart';

import '../models/home_device.dart';
import '../theme/hata_colors.dart';
import '../widgets/device_card.dart';
import '../widgets/hata_top_bar.dart';
import '../widgets/home_fab_menu.dart';
import 'shopping_screen.dart';

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
      floatingActionButton: HomeFabMenu(
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
