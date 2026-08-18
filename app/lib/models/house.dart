class House {
  const House({
    required this.id,
    required this.displayName,
    required this.role,
    this.location,
    this.weather,
  });

  final String id;
  final String displayName;
  final String role;
  final String? location;
  final WeatherInfo? weather;

  factory House.fromJson(Map<String, dynamic> json) {
    final id = json['id']?.toString() ?? '';
    final weatherJson = json['weather'];
    return House(
      id: id,
      displayName: json['displayName']?.toString() ?? id,
      role: json['role']?.toString() ?? '',
      location: json['location']?.toString(),
      weather: weatherJson is Map<String, dynamic>
          ? WeatherInfo.fromJson(weatherJson)
          : null,
    );
  }
}

class WeatherInfo {
  const WeatherInfo({
    required this.status,
    this.locationLabel,
    this.temperature,
    this.temperatureUnit,
    this.conditionCode,
    this.conditionText,
    this.conditionIcon,
    this.humidityPercent,
    this.windSpeed,
    this.windSpeedUnit,
    this.observedAt,
    this.updatedAt,
    this.error,
  });

  final String status;
  final String? locationLabel;
  final double? temperature;
  final String? temperatureUnit;
  final int? conditionCode;
  final String? conditionText;
  final String? conditionIcon;
  final double? humidityPercent;
  final double? windSpeed;
  final String? windSpeedUnit;
  final DateTime? observedAt;
  final DateTime? updatedAt;
  final String? error;

  bool get hasCurrentConditions =>
      (status == 'ok' || status == 'stale') && temperature != null;

  factory WeatherInfo.fromJson(Map<String, dynamic> json) {
    return WeatherInfo(
      status: json['status']?.toString() ?? '',
      locationLabel: _stringOrNull(json['locationLabel']),
      temperature: _doubleOrNull(json['temperature']),
      temperatureUnit: _stringOrNull(json['temperatureUnit']),
      conditionCode: _intOrNull(json['conditionCode']),
      conditionText: _stringOrNull(json['conditionText']),
      conditionIcon: _stringOrNull(json['conditionIcon']),
      humidityPercent: _doubleOrNull(json['humidityPercent']),
      windSpeed: _doubleOrNull(json['windSpeed']),
      windSpeedUnit: _stringOrNull(json['windSpeedUnit']),
      observedAt: _dateTimeOrNull(json['observedAt']),
      updatedAt: _dateTimeOrNull(json['updatedAt']),
      error: _stringOrNull(json['error']),
    );
  }

  static String? _stringOrNull(Object? value) {
    final text = value?.toString().trim();
    return text == null || text.isEmpty ? null : text;
  }

  static double? _doubleOrNull(Object? value) {
    if (value is num) return value.toDouble();
    return double.tryParse(value?.toString() ?? '');
  }

  static int? _intOrNull(Object? value) {
    if (value is int) return value;
    return int.tryParse(value?.toString() ?? '');
  }

  static DateTime? _dateTimeOrNull(Object? value) =>
      DateTime.tryParse(value?.toString() ?? '');
}
