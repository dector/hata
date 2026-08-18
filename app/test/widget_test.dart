import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:app/main.dart';

Future<void> _login(WidgetTester tester) async {
  await tester.pumpWidget(const HataApp());

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
    await tester.pumpWidget(const HataApp());

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
