import 'package:app/api/api.dart';
import 'package:app/models/house.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';

void main() {
  test('House parses weather from extension extras', () {
    final house = House.fromJson({
      'id': 'home',
      'displayName': 'My Home',
      'role': 'owner',
      'extras': {
        weatherExtensionId: {
          'status': 'ok',
          'temperature': 21.5,
          'conditionText': 'Clear',
        },
      },
    });

    expect(house.weather?.status, 'ok');
    expect(house.weather?.temperature, 21.5);
    expect(house.weather?.conditionText, 'Clear');
  });

  test('House still parses legacy weather field', () {
    final house = House.fromJson({
      'id': 'home',
      'displayName': 'My Home',
      'role': 'owner',
      'weather': {'status': 'pending'},
    });

    expect(house.weather?.status, 'pending');
  });

  test('refreshWeather calls generic house extension action endpoint', () async {
    final api = Api(
      httpClient: MockClient((request) async {
        expect(request.method, 'POST');
        expect(
          request.url.path,
          '/api/latest/house/home/extension/$weatherExtensionId/actions/refresh',
        );
        expect(request.headers['Authorization'], 'Bearer token');
        return http.Response('{"weather":{"status":"ok"}}', 200);
      }),
    );

    final weather = await api.refreshWeather(
      'https://hata.test',
      'token',
      'home',
    );

    expect(weather?.status, 'ok');
  });
}
