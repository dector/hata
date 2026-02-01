package hata.data.repositories

import android.content.SharedPreferences
import androidx.core.content.edit
import hata.data.models.Session
import hata.data.models.SessionUser
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow


internal class SessionRepositoryImpl(
    private val sharedPreferences: SharedPreferences,
) : SessionRepository {

    private val sessionFlow = MutableStateFlow(loadSession())

    override fun observeSession(): Flow<Session?> {
        return sessionFlow.asStateFlow()
    }

    override suspend fun getSession(): Session? {
        return sessionFlow.value
    }

    override suspend fun saveSession(session: Session) {
        sharedPreferences.edit().apply {
            putString(KEY_TOKEN, session.token)
            putString(KEY_USER_NAME, session.user.name)
            putBoolean(KEY_IS_ACTIVE, session.isActive)
            putString(KEY_SERVER_URL, session.serverUrl)
            apply()
        }
        sessionFlow.value = session
    }

    override suspend fun clearSession() {
        sharedPreferences.edit { clear() }
        sessionFlow.value = null
    }

    private fun loadSession(): Session? {
        val token = sharedPreferences.getString(KEY_TOKEN, null)
        val userName = sharedPreferences.getString(KEY_USER_NAME, null)
        val isActive = sharedPreferences.getBoolean(KEY_IS_ACTIVE, false)
        val serverUrl = sharedPreferences.getString(KEY_SERVER_URL, null)

        return if (token != null && userName != null && serverUrl != null) {
            Session(
                token = token,
                user = SessionUser(name = userName),
                isActive = isActive,
                serverUrl = serverUrl,
            )
        } else {
            null
        }
    }

    private companion object {
        const val KEY_TOKEN = "session_token"
        const val KEY_USER_NAME = "session_user_name"
        const val KEY_IS_ACTIVE = "session_is_active"
        const val KEY_SERVER_URL = "session_server_url"
    }
}
