import '../api/api.dart';
import '../models/house.dart';
import '../session/session_database.dart';
import '../session/session_repository.dart';

class WearWeatherBridge {
  WearWeatherBridge({Api? api, SessionRepository? sessionRepository})
    : _api = api ?? Api(),
      _sessionRepository = sessionRepository ?? SqliteSessionRepository();

  static final WearWeatherBridge instance = WearWeatherBridge();

  final Api _api;
  final SessionRepository _sessionRepository;
  WeatherInfo? _latestWeather;

  void updateFromHouses(List<House> houses) {
    _latestWeather = _firstCurrentWeather(houses) ?? _latestWeather;
  }

  void updateWeather(WeatherInfo weather) {
    if (weather.hasCurrentConditions) {
      _latestWeather = weather;
    }
  }

  Future<Map<String, Object>> getWeather() async {
    final cached = _latestWeather;
    if (cached != null && cached.hasCurrentConditions) {
      return _payload(cached);
    }

    final session = await _sessionRepository.load();
    if (session == null || session.isExpired) {
      throw Exception('Open Hata and sign in on your phone');
    }

    final houses = await _api.fetchHouses(session.serverUrl, session.token);
    final weather = _firstCurrentWeather(houses);
    if (weather == null) {
      throw Exception('Weather unavailable');
    }

    _latestWeather = weather;
    return _payload(weather);
  }

  WeatherInfo? _firstCurrentWeather(List<House> houses) {
    for (final house in houses) {
      final weather = house.weather;
      if (weather != null && weather.hasCurrentConditions) {
        return weather;
      }
    }
    return null;
  }

  Map<String, Object> _payload(WeatherInfo weather) {
    final temperature = weather.temperature;
    if (temperature == null) {
      throw Exception('Weather temperature unavailable');
    }

    return <String, Object>{
      'temperatureC': temperature.round(),
      'temperatureLabel': _temperatureLabel(weather),
      'condition': weather.conditionText ?? 'Weather now',
      'updatedAt': (weather.updatedAt ?? weather.observedAt ?? DateTime.now())
          .toUtc()
          .toIso8601String(),
    };
  }

  String _temperatureLabel(WeatherInfo weather) {
    final value = weather.temperature ?? 0;
    final text = value == value.roundToDouble()
        ? value.toStringAsFixed(0)
        : value.toStringAsFixed(1);
    final unit = weather.temperatureUnit?.trim();
    if (unit == null || unit.isEmpty) return text;
    return unit.startsWith('°') ? '$text$unit' : '$text °$unit';
  }
}
