import 'package:flutter/material.dart';

import 'api/hata_api_client.dart';
import 'app_shell.dart';
import 'routes.dart';
import 'session/session_repository.dart';
import 'theme/hata_colors.dart';

class HataApp extends StatelessWidget {
  const HataApp({super.key, this.apiClient, this.sessionRepository});

  final HataApiClient? apiClient;
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
          apiClient: apiClient,
          sessionRepository: sessionRepository,
        ),
      },
    );
  }
}
