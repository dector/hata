import 'package:flutter/material.dart';

class HomeDevice {
  const HomeDevice({
    required this.id,
    required this.title,
    required this.status,
    required this.icon,
    required this.isOn,
  });

  final String id;
  final String title;
  final String status;
  final IconData icon;
  final bool isOn;

  HomeDevice copyWith({bool? isOn}) {
    return HomeDevice(
      id: id,
      title: title,
      status: status,
      icon: icon,
      isOn: isOn ?? this.isOn,
    );
  }
}
