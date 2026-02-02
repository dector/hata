package hata.di

import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import hata.BuildConfig
import hata.feature.profile.model.AppInfo
import hata.feature.profile.model.ProfileSession
import hata.feature.profile.repository.ProfileSessionRepository
import hata.feature.session.repository.SessionRepository
import javax.inject.Singleton


@Module
@InstallIn(SingletonComponent::class)
object ProfileFeatureModule {

    @Provides
    @Singleton
    fun profileSessionRepository(
        sessionRepository: SessionRepository,
    ): ProfileSessionRepository {
        return RealProfileSessionRepository(sessionRepository)
    }

    @Provides
    fun appInfo(): AppInfo = AppInfo(
        version = BuildConfig.VERSION_NAME,
    )
}

class RealProfileSessionRepository(
    private val sessionRepository: SessionRepository,
) : ProfileSessionRepository {

    override suspend fun getSession(): ProfileSession? {
        val session = sessionRepository.getSession() ?: return null

        return ProfileSession(
            userName = session.user.name,
            isActive = session.isActive,
            serverUrl = session.serverUrl,
        )
    }

    override suspend fun clearSession() {
        sessionRepository.clearSession()
    }
}
