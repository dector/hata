import 'dart:convert';

import 'package:http/http.dart' as http;

import '../models/api_device.dart';
import '../models/auth_session.dart';
import '../models/house.dart';
import '../models/server_info.dart';
import 'api_error.dart';

class HataApiClient {
  HataApiClient({http.Client? httpClient})
    : _httpClient = httpClient ?? http.Client();

  final http.Client _httpClient;

  String normalizeBaseUrl(String serverUrl) =>
      serverUrl.trim().replaceFirst(RegExp(r'/+$'), '');

  Future<ServerInfo> ping(String serverUrl) async {
    final response = await _httpClient.get(_uri(serverUrl, '/api/latest/ping'));
    final json = _decodeResponse(response);
    return ServerInfo.fromJson(json);
  }

  Future<AuthSession> login(
    String serverUrl,
    String username,
    String password,
  ) async {
    final response = await _httpClient.post(
      _uri(serverUrl, '/api/latest/auth/login'),
      headers: const {'Content-Type': 'application/json'},
      body: jsonEncode({'username': username, 'password': password}),
    );
    final json = _decodeResponse(response);
    final session = AuthSession.fromJson(json);
    if (session.token.isEmpty) {
      throw const ApiError('Server returned an empty session token');
    }
    return session;
  }

  Future<void> fetchHouse(String serverUrl, String token) async {
    await fetchHouses(serverUrl, token);
  }

  Future<List<House>> fetchHouses(String serverUrl, String token) async {
    final response = await _httpClient.get(
      _uri(serverUrl, '/api/latest/house'),
      headers: _authHeaders(token),
    );
    final json = _decodeResponse(response);
    final houses = json['houses'] as List<dynamic>? ?? const [];
    return houses
        .whereType<Map<String, dynamic>>()
        .map(House.fromJson)
        .where((house) => house.id.isNotEmpty)
        .toList();
  }

  Future<List<ApiDevice>> fetchDevices(String serverUrl, String token) async {
    final response = await _httpClient.get(
      _uri(serverUrl, '/api/latest/device'),
      headers: _authHeaders(token),
    );
    final json = _decodeResponse(response);
    return _decodeDevices(json);
  }

  Future<List<ApiDevice>> fetchHouseDevices(
    String serverUrl,
    String token,
    String houseId,
  ) async {
    final response = await _httpClient.get(
      _uri(serverUrl, '/api/latest/house/$houseId/device'),
      headers: _authHeaders(token),
    );
    final json = _decodeResponse(response);
    return _decodeDevices(json);
  }

  Future<DeviceStateResult> setDeviceState(
    String serverUrl,
    String token,
    String houseId,
    String deviceId,
    bool isOn,
  ) async {
    final response = await _httpClient.patch(
      _uri(serverUrl, '/api/latest/house/$houseId/device/$deviceId/state'),
      headers: _authHeaders(token, json: true),
      body: jsonEncode({'state': isOn ? 'on' : 'off'}),
    );
    final json = _decodeResponse(response);
    return DeviceStateResult.fromJson(json);
  }

  List<ApiDevice> _decodeDevices(Map<String, dynamic> json) {
    final devices = json['devices'] as List<dynamic>? ?? const [];
    return devices
        .whereType<Map<String, dynamic>>()
        .map(ApiDevice.fromJson)
        .where((device) => device.id.isNotEmpty)
        .toList();
  }

  Map<String, String> _authHeaders(String token, {bool json = false}) => {
    'Authorization': 'Bearer $token',
    if (json) 'Content-Type': 'application/json',
  };

  Uri _uri(String serverUrl, String path) =>
      Uri.parse('${normalizeBaseUrl(serverUrl)}$path');

  Map<String, dynamic> _decodeResponse(http.Response response) {
    Map<String, dynamic> json;
    try {
      json = response.body.isEmpty
          ? <String, dynamic>{}
          : jsonDecode(response.body) as Map<String, dynamic>;
    } catch (_) {
      json = <String, dynamic>{};
    }

    if (response.statusCode < 200 || response.statusCode >= 300) {
      final error = json['error'];
      if (error is Map<String, dynamic>) {
        throw ApiError(
          error['message'] as String? ?? 'Request failed',
          code: error['code'] as String?,
          statusCode: response.statusCode,
        );
      }
      throw ApiError(
        'Request failed (${response.statusCode})',
        statusCode: response.statusCode,
      );
    }

    return json;
  }
}
