import 'package:app/api/api.dart';
import 'package:app/app.dart';
import 'package:app/models/api_device.dart';
import 'package:app/models/auth_session.dart';
import 'package:app/models/house.dart';
import 'package:app/models/server_info.dart';
import 'package:app/models/shopping_item.dart';
import 'package:app/session/session.dart';
import 'package:app/session/session_repository.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

class FakeApiClient extends Api {
  @override
  Future<ServerInfo> ping(String serverUrl) async => const ServerInfo(
    serverName: 'Hata Server',
    version: '1.0.0',
    apiVersion: 'latest',
  );

  @override
  Future<AuthSession> login(
    String serverUrl,
    String username,
    String password,
  ) async => AuthSession(
    token: 'token',
    validUntil: DateTime.now().add(const Duration(days: 1)),
    displayName: 'Dan',
  );

  @override
  Future<void> fetchHouse(String serverUrl, String token) async {}

  @override
  Future<List<House>> fetchHouses(String serverUrl, String token) async =>
      const [House(id: 'home', displayName: 'My Home', role: 'owner')];

  @override
  Future<List<ApiDevice>> fetchDevices(String serverUrl, String token) async =>
      const [
        ApiDevice(
          id: 'cooler',
          name: 'Air Cooler',
          integration: DeviceIntegration(id: 'test', data: {}),
          state: 'on',
          availability: 'online',
          capabilities: DeviceCapabilities(
            light: false,
            brightness: false,
            colorPresets: false,
          ),
          houseId: 'home',
        ),
        ApiDevice(
          id: 'kitchen',
          name: 'Kitchen Light',
          integration: DeviceIntegration(id: 'test-light', data: {}),
          state: 'off',
          availability: 'online',
          capabilities: DeviceCapabilities(
            light: true,
            brightness: false,
            colorPresets: false,
          ),
          houseId: 'home',
        ),
        ApiDevice(
          id: 'office',
          name: 'Office Light',
          integration: DeviceIntegration(id: 'test-light', data: {}),
          state: 'on',
          availability: 'offline',
          capabilities: DeviceCapabilities(
            light: true,
            brightness: false,
            colorPresets: false,
          ),
          houseId: 'home',
        ),
      ];

  @override
  Future<List<ShoppingItem>> fetchShoppingItems(
    String serverUrl,
    String token,
    String houseId,
    String listId,
  ) async => const [
    ShoppingItem(
      uid: 'coffee',
      name: 'Coffee beans',
      position: 1,
      checkedAt: null,
      checkedByUserId: null,
      deletedAt: null,
      createdAt: null,
      updatedAt: null,
    ),
  ];
}

class MemorySessionRepository implements SessionRepository {
  Session? session;

  @override
  Future<void> clear() async => session = null;

  @override
  Future<Session?> load() async => session;

  @override
  Future<void> save(Session session) async => this.session = session;
}

Future<void> _login(WidgetTester tester) async {
  await tester.pumpWidget(
    HataApp(
      apiClient: FakeApiClient(),
      sessionRepository: MemorySessionRepository(),
    ),
  );
  await tester.pumpAndSettle();

  expect(
    find.widgetWithText(TextField, 'http://10.0.2.2:4501'),
    findsOneWidget,
  );
  await tester.tap(find.text('Connect'));
  await tester.pumpAndSettle();

  await tester.enterText(find.byType(TextField).at(0), 'dan');
  await tester.enterText(find.byType(TextField).at(1), 'secret');
  await tester.pump();
  await tester.tap(find.text('Login'));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('renders login from app shell', (WidgetTester tester) async {
    await tester.pumpWidget(
      HataApp(
        apiClient: FakeApiClient(),
        sessionRepository: MemorySessionRepository(),
      ),
    );
    await tester.pumpAndSettle();

    expect(find.text('Hata'), findsOneWidget);
    expect(find.text('Server URL'), findsOneWidget);
    expect(find.text('Connect'), findsOneWidget);
  });

  testWidgets('renders Hata home dashboard after login', (
    WidgetTester tester,
  ) async {
    await _login(tester);

    expect(find.text('My Home'), findsOneWidget);
    expect(find.text('Welcome home,\nDan'), findsOneWidget);
    expect(find.text('Air Cooler'), findsOneWidget);
    expect(find.text('Kitchen Light'), findsOneWidget);
    expect(find.text('Office Light'), findsOneWidget);
    expect(find.text('test-light • Offline'), findsOneWidget);

    final switches = tester.widgetList<Switch>(find.byType(Switch)).toList();
    expect(switches[2].value, isFalse);
    expect(switches[2].onChanged, isNull);
  });

  testWidgets('opens shopping screen from floating menu', (
    WidgetTester tester,
  ) async {
    await _login(tester);

    await tester.tap(find.byType(FloatingActionButton));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Shopping'));
    await tester.pumpAndSettle();

    expect(find.text('Shopping'), findsOneWidget);
    expect(find.text('Coffee beans'), findsOneWidget);
  });
}
