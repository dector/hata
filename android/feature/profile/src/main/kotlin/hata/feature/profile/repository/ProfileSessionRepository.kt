package hata.feature.profile.repository

import hata.feature.profile.model.ProfileSession


interface ProfileSessionRepository {
    suspend fun getSession(): ProfileSession?
    suspend fun clearSession()
}
