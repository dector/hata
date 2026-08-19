import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'app.dart';
import 'wear/wear_weather_bridge.dart';

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
        return WearWeatherBridge.instance.getWeather();
      default:
        throw PlatformException(
          code: 'not_implemented',
          message: 'Unknown wear method: ${call.method}',
        );
    }
  });
}
