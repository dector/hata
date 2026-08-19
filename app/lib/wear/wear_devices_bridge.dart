import '../api/api.dart';
import '../models/api_device.dart';
import '../session/session.dart';
import '../session/session_database.dart';
import '../session/session_repository.dart';

class WearDevicesBridge {
  WearDevicesBridge({Api? api, SessionRepository? sessionRepository})
    : _api = api ?? Api(),
      _sessionRepository = sessionRepository ?? SqliteSessionRepository();

  static final WearDevicesBridge instance = WearDevicesBridge();

  final Api _api;
  final SessionRepository _sessionRepository;
  List<ApiDevice>? _latestDevices;

  void updateDevices(List<ApiDevice> devices) {
    _latestDevices = devices;
  }

  Future<List<Map<String, Object?>>> getDevices() async {
    final devices = await _loadDevices();
    return devices.map(_payload).toList();
  }

  Future<List<Map<String, Object?>>> toggleDevice(Map<Object?, Object?> args) async {
    final houseId = args['houseId']?.toString() ?? '';
    final deviceId = args['deviceId']?.toString() ?? '';
    final targetState = args['targetState']?.toString().toLowerCase() ?? '';

    if (houseId.isEmpty || deviceId.isEmpty) {
      throw Exception('Missing device identifiers');
    }
    if (targetState != 'on' && targetState != 'off') {
      throw Exception('Unsupported target state');
    }

    final session = await _validSession();
    await _api.setDeviceState(
      session.serverUrl,
      session.token,
      houseId,
      deviceId,
      targetState == 'on',
    );

    final devices = await _api.fetchDevices(session.serverUrl, session.token);
    _latestDevices = devices;
    return devices.map(_payload).toList();
  }

  Future<List<ApiDevice>> _loadDevices() async {
    final session = await _validSession();
    final devices = await _api.fetchDevices(session.serverUrl, session.token);
    _latestDevices = devices;
    return devices;
  }

  Future<Session> _validSession() async {
    final session = await _sessionRepository.load();
    if (session == null || session.isExpired) {
      throw Exception('Open Hata and sign in on your phone');
    }
    return session;
  }

  Map<String, Object?> _payload(ApiDevice device) {
    return <String, Object?>{
      'id': device.id,
      'houseId': device.houseId,
      'name': device.name.isEmpty ? device.id : device.name,
      'state': device.state,
      'availability': device.availability,
      'brightness': _brightness(device),
    };
  }

  int? _brightness(ApiDevice device) {
    final value = device.light?['brightness'];
    if (value is int) return value.clamp(0, 100).toInt();
    if (value is num) return value.round().clamp(0, 100).toInt();
    return int.tryParse(value?.toString() ?? '')?.clamp(0, 100).toInt();
  }
}
