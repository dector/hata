import 'package:flutter/material.dart';

import 'api/api.dart';
import 'app_shell.dart';
import 'routes.dart';
import 'session/session_repository.dart';
import 'theme/hata_colors.dart';

class HataApp extends StatelessWidget {
  const HataApp({super.key, this.apiClient, this.sessionRepository});

  final Api? apiClient;
  final SessionRepository? sessionRepository;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Hata',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: HataColors.primary,
          brightness: Brightness.dark,
        ),
        scaffoldBackgroundColor: HataColors.background,
        useMaterial3: true,
      ),
      initialRoute: AppRoute.init,
      routes: {
        AppRoute.init: (_) => AppShell(
          api: apiClient,
          sessionRepository: sessionRepository,
        ),
      },
    );
  }
}
