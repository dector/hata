import 'package:path/path.dart' as p;
import 'package:sqflite/sqflite.dart';

import 'session.dart';
import 'session_repository.dart';

class SqliteSessionRepository implements SessionRepository {
  Database? _database;

  Future<Database> get _db async {
    final existing = _database;
    if (existing != null) return existing;

    final dbPath = p.join(await getDatabasesPath(), 'hata_session.db');
    return _database = await openDatabase(
      dbPath,
      version: 1,
      onCreate: (db, version) async {
        await db.execute('''
CREATE TABLE session (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  token TEXT NOT NULL,
  server_url TEXT NOT NULL,
  username TEXT NOT NULL,
  display_name TEXT,
  valid_until TEXT NOT NULL,
  created_at TEXT NOT NULL
)
''');
      },
    );
  }

  @override
  Future<Session?> load() async {
    final rows = await (await _db).query('session', where: 'id = 1', limit: 1);
    if (rows.isEmpty) return null;
    final row = rows.single;
    return Session(
      token: row['token'] as String,
      serverUrl: row['server_url'] as String,
      username: row['username'] as String,
      displayName: row['display_name'] as String?,
      validUntil: DateTime.parse(row['valid_until'] as String),
      createdAt: DateTime.parse(row['created_at'] as String),
    );
  }

  @override
  Future<void> save(Session session) async {
    await (await _db).insert('session', {
      'id': 1,
      'token': session.token,
      'server_url': session.serverUrl,
      'username': session.username,
      'display_name': session.displayName,
      'valid_until': session.validUntil.toIso8601String(),
      'created_at': session.createdAt.toIso8601String(),
    }, conflictAlgorithm: ConflictAlgorithm.replace);
  }

  @override
  Future<void> clear() async {
    await (await _db).delete('session');
  }
}
