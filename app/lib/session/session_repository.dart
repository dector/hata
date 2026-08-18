import 'session.dart';

abstract class SessionRepository {
  Future<Session?> load();
  Future<void> save(Session session);
  Future<void> clear();
}
