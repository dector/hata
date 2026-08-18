import 'dart:convert';

import 'package:http/http.dart' as http;

import '../models/auth_session.dart';
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
    final response = await _httpClient.get(
      _uri(serverUrl, '/api/latest/house'),
      headers: {'Authorization': 'Bearer $token'},
    );
    _decodeResponse(response);
  }

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
