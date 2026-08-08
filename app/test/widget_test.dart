import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:app/main.dart';

void main() {
  testWidgets('renders Hata home dashboard', (WidgetTester tester) async {
    await tester.pumpWidget(const HataApp());

    expect(find.text('My Home'), findsOneWidget);
    expect(find.text('Welcome home,\nDan'), findsOneWidget);
    expect(find.text('Air Cooler'), findsOneWidget);
    expect(find.text('Kitchen Light'), findsOneWidget);
  });

  testWidgets('opens shopping screen from floating menu', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(const HataApp());

    await tester.tap(find.byType(FloatingActionButton));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Shopping'));
    await tester.pumpAndSettle();

    expect(find.text('Shopping'), findsOneWidget);
    expect(find.text('Coffee beans'), findsOneWidget);
  });
}
