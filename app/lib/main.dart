import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'app.dart';
import 'wear/wear_devices_bridge.dart';
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
      case 'getDevices':
        return WearDevicesBridge.instance.getDevices();
      case 'toggleDevice':
        final args = call.arguments;
        if (args is Map<Object?, Object?>) {
          return WearDevicesBridge.instance.toggleDevice(args);
        }
        throw PlatformException(
          code: 'invalid_arguments',
          message: 'toggleDevice expects a map payload',
        );
      default:
        throw PlatformException(
          code: 'not_implemented',
          message: 'Unknown wear method: ${call.method}',
        );
    }
  });
}
