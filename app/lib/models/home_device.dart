import 'package:flutter/material.dart';

class HomeDevice {
  const HomeDevice({
    required this.id,
    required this.houseId,
    required this.title,
    required this.status,
    required this.icon,
    required this.isOn,
    required this.canControl,
    required this.rawState,
    required this.availability,
    this.isAwaitingConfirmation = false,
  });

  final String id;
  final String houseId;
  final String title;
  final String status;
  final IconData icon;
  final bool isOn;
  final bool canControl;
  final bool isAwaitingConfirmation;
  final String rawState;
  final String availability;

  HomeDevice copyWith({
    String? status,
    bool? isOn,
    bool? canControl,
    bool? isAwaitingConfirmation,
    String? rawState,
    String? availability,
  }) {
    return HomeDevice(
      id: id,
      houseId: houseId,
      title: title,
      status: status ?? this.status,
      icon: icon,
      isOn: isOn ?? this.isOn,
      canControl: canControl ?? this.canControl,
      isAwaitingConfirmation:
          isAwaitingConfirmation ?? this.isAwaitingConfirmation,
      rawState: rawState ?? this.rawState,
      availability: availability ?? this.availability,
    );
  }
}
