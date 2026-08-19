import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'app.dart';

export 'app.dart';

const _wearChannel = MethodChannel('hata/wear');

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  _registerWearBridge();
  runApp(const HataApp());
}

void _registerWearBridge() {
  _wearChannel.setMethodCallHandler((call) async {
    switch (call.method) {
      case 'getWeather':
        return <String, Object>{
          'temperatureC': 22,
          'condition': 'Cloudy',
          'updatedAt': DateTime.now().toUtc().toIso8601String(),
        };
      default:
        throw PlatformException(
          code: 'not_implemented',
          message: 'Unknown wear method: ${call.method}',
        );
    }
  });
}
