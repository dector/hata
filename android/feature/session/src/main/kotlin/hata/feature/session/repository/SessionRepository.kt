package hata.feature.session.repository

import hata.feature.session.model.Session
import kotlinx.coroutines.flow.Flow


interface SessionRepository {
    fun observeSession(): Flow<Session?>
    suspend fun getSession(): Session?
    suspend fun saveSession(session: Session)
    suspend fun clearSession()
}
