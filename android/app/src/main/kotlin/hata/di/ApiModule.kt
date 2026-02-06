package hata.di

import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import hata.data.api.AuthTokenProvider
import hata.data.api.RealServerProbeApi
import hata.data.api.RealServerServiceFactory
import hata.data.api.ServerProbeApi
import hata.data.api.ServerServiceFactory
import hata.feature.session.repository.SessionRepository
import javax.inject.Singleton


@Module
@InstallIn(SingletonComponent::class)
object ApiModule {

    @Provides
    @Singleton
    fun serverProbeApi(): ServerProbeApi {
        return RealServerProbeApi()
    }

    @Provides
    @Singleton
    fun serverServiceFactory(
        authTokenProvider: AuthTokenProvider,
    ): ServerServiceFactory {
        return RealServerServiceFactory(authTokenProvider)
    }

    @Provides
    @Singleton
    fun authTokenProvider(
        sessionRepository: SessionRepository,
    ): AuthTokenProvider {
        return RealAuthTokenProvider(sessionRepository)
    }
}

internal class RealAuthTokenProvider(
    private val sessionRepository: SessionRepository,
) : AuthTokenProvider {

    override fun getToken(): String? =
        sessionRepository.currentSession()?.token
}
