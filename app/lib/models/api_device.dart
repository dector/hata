class DeviceIntegration {
  const DeviceIntegration({required this.id, required this.data});

  final String id;
  final Map<String, dynamic> data;

  factory DeviceIntegration.fromJson(Map<String, dynamic>? json) {
    return DeviceIntegration(
      id: json?['id']?.toString() ?? 'unknown',
      data: json?['data'] is Map<String, dynamic>
          ? json!['data'] as Map<String, dynamic>
          : <String, dynamic>{},
    );
  }
}

class DeviceCapabilities {
  const DeviceCapabilities({
    required this.light,
    required this.brightness,
    required this.colorPresets,
  });

  final bool light;
  final bool brightness;
  final bool colorPresets;

  factory DeviceCapabilities.fromJson(Map<String, dynamic>? json) {
    return DeviceCapabilities(
      light: json?['light'] == true,
      brightness: json?['brightness'] == true,
      colorPresets: json?['colorPresets'] == true,
    );
  }
}

class ApiDevice {
  const ApiDevice({
    required this.id,
    required this.name,
    required this.integration,
    required this.state,
    required this.availability,
    required this.capabilities,
    required this.houseId,
    this.light,
  });

  final String id;
  final String name;
  final DeviceIntegration integration;
  final String state;
  final String availability;
  final DeviceCapabilities capabilities;
  final String houseId;
  final Map<String, dynamic>? light;

  factory ApiDevice.fromJson(Map<String, dynamic> json) {
    final house = json['house'] is Map<String, dynamic>
        ? json['house'] as Map<String, dynamic>
        : const <String, dynamic>{};

    return ApiDevice(
      id: json['id']?.toString() ?? '',
      name: json['name']?.toString() ?? '',
      integration: DeviceIntegration.fromJson(
        json['integration'] is Map<String, dynamic>
            ? json['integration'] as Map<String, dynamic>
            : null,
      ),
      state: json['state']?.toString() ?? 'unknown',
      availability: json['availability']?.toString() ?? 'unknown',
      capabilities: DeviceCapabilities.fromJson(
        json['capabilities'] is Map<String, dynamic>
            ? json['capabilities'] as Map<String, dynamic>
            : null,
      ),
      houseId: house['id']?.toString() ?? '',
      light: json['light'] is Map<String, dynamic>
          ? json['light'] as Map<String, dynamic>
          : null,
    );
  }
}

class DeviceStateResult {
  const DeviceStateResult({
    required this.deviceId,
    required this.houseId,
    required this.state,
  });

  final String deviceId;
  final String houseId;
  final String state;

  factory DeviceStateResult.fromJson(Map<String, dynamic> json) {
    final house = json['house'] is Map<String, dynamic>
        ? json['house'] as Map<String, dynamic>
        : const <String, dynamic>{};

    return DeviceStateResult(
      deviceId: json['deviceId']?.toString() ?? '',
      houseId: house['id']?.toString() ?? '',
      state: json['state']?.toString() ?? 'unknown',
    );
  }
}
