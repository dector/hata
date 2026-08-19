import 'dart:convert';

import 'package:http/http.dart' as http;

import '../models/api_device.dart';
import '../models/auth_session.dart';
import '../models/house.dart';
import '../models/server_info.dart';
import '../models/shopping_item.dart';
import '../models/shopping_list.dart';
import 'api_error.dart';

const String apiUrlBase = '/api/latest';

String _apiUrl(String path) {
  final normalizedPath = path.startsWith('/') ? path.substring(1) : path;
  return '$apiUrlBase/$normalizedPath';
}

class Api {
  Api({http.Client? httpClient}) : _client = httpClient ?? http.Client();

  final http.Client _client;

  String normalizeBaseUrl(String serverUrl) =>
      serverUrl.trim().replaceFirst(RegExp(r'/+$'), '');

  Future<ServerInfo> ping(String serverUrl) async {
    final response = await _client.get(_uri(serverUrl, _apiUrl('ping')));
    final json = _decodeResponse(response);
    return ServerInfo.fromJson(json);
  }

  Future<AuthSession> login(
    String serverUrl,
    String username,
    String password,
  ) async {
    final response = await _client.post(
      _uri(serverUrl, _apiUrl('auth/login')),
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
    final response = await _client.get(
      _uri(serverUrl, _apiUrl('house')),
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
    final response = await _client.get(
      _uri(serverUrl, _apiUrl('device')),
      headers: _authHeaders(token),
    );
    final json = _decodeResponse(response);
    return _decodeDevices(json);
  }

  Future<WeatherInfo?> refreshWeather(
    String serverUrl,
    String token,
    String houseId,
  ) async {
    final response = await _client.post(
      _uri(
        serverUrl,
        _apiUrl(
          'house/$houseId/extension/${Uri.encodeComponent(weatherExtensionId)}/actions/refresh',
        ),
      ),
      headers: _authHeaders(token),
    );
    final json = _decodeResponse(response);
    final weather = json['weather'];
    return weather is Map<String, dynamic>
        ? WeatherInfo.fromJson(weather)
        : null;
  }

  Future<List<ApiDevice>> fetchHouseDevices(
    String serverUrl,
    String token,
    String houseId,
  ) async {
    final response = await _client.get(
      _uri(serverUrl, _apiUrl('house/$houseId/device')),
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
    final response = await _client.patch(
      _uri(serverUrl, _apiUrl('house/$houseId/device/$deviceId/state')),
      headers: _authHeaders(token, json: true),
      body: jsonEncode({'state': isOn ? 'on' : 'off'}),
    );
    final json = _decodeResponse(response);
    return DeviceStateResult.fromJson(json);
  }

  Future<List<ShoppingList>> fetchShoppingLists(
    String serverUrl,
    String token,
    String houseId,
  ) async {
    final response = await _client.get(
      _uri(serverUrl, _apiUrl('house/$houseId/shopping-list')),
      headers: _authHeaders(token),
    );
    final json = _decodeResponse(response);
    final lists = json['shoppingLists'] as List<dynamic>? ?? const [];
    return lists
        .whereType<Map<String, dynamic>>()
        .map(ShoppingList.fromJson)
        .where((list) => list.uid.isNotEmpty)
        .toList();
  }

  Future<List<ShoppingItem>> fetchShoppingItems(
    String serverUrl,
    String token,
    String houseId,
    String listId,
  ) async {
    final response = await _client.get(
      _uri(serverUrl, _apiUrl('house/$houseId/shopping-list/$listId/item')),
      headers: _authHeaders(token),
    );
    final json = _decodeResponse(response);
    final items = json['items'] as List<dynamic>? ?? const [];
    return _sortShoppingItems(
      items
          .whereType<Map<String, dynamic>>()
          .map(ShoppingItem.fromJson)
          .where((item) => item.uid.isNotEmpty)
          .toList(),
    );
  }

  Future<ShoppingItem> createShoppingItem(
    String serverUrl,
    String token,
    String houseId,
    String listId,
    String name,
  ) async {
    final response = await _client.post(
      _uri(serverUrl, _apiUrl('house/$houseId/shopping-list/$listId/item')),
      headers: _authHeaders(token, json: true),
      body: jsonEncode({'name': name}),
    );
    final json = _decodeResponse(response);
    return ShoppingItem.fromJson(json['item'] as Map<String, dynamic>? ?? {});
  }

  Future<ShoppingItem> setShoppingItemChecked(
    String serverUrl,
    String token,
    String houseId,
    String listId,
    String itemId,
    bool checked,
  ) async {
    final response = await _client.patch(
      _uri(
        serverUrl,
        _apiUrl('house/$houseId/shopping-list/$listId/item/$itemId/check'),
      ),
      headers: _authHeaders(token, json: true),
      body: jsonEncode({'checked': checked}),
    );
    final json = _decodeResponse(response);
    return ShoppingItem.fromJson(json['item'] as Map<String, dynamic>? ?? {});
  }

  List<ApiDevice> _decodeDevices(Map<String, dynamic> json) {
    final devices = json['devices'] as List<dynamic>? ?? const [];
    return devices
        .whereType<Map<String, dynamic>>()
        .map(ApiDevice.fromJson)
        .where((device) => device.id.isNotEmpty)
        .toList();
  }

  List<ShoppingItem> _sortShoppingItems(List<ShoppingItem> items) {
    return List<ShoppingItem>.from(items)..sort((a, b) {
      final checkedComparison = a.isChecked == b.isChecked
          ? 0
          : a.isChecked
          ? 1
          : -1;
      if (checkedComparison != 0) return checkedComparison;

      final nameComparison = a.name.toLowerCase().compareTo(
        b.name.toLowerCase(),
      );
      if (nameComparison != 0) return nameComparison;

      return a.uid.compareTo(b.uid);
    });
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
