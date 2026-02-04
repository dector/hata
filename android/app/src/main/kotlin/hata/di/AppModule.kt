package hata.di

import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import hata.core.annotations.ComputeDispatcher
import hata.core.annotations.IoDispatcher
import hata.data.api.RealServerServiceImpl
import hata.data.api.RemoteConfigurationApi
import hata.data.api.RemoteConfigurationServiceImpl
import hata.data.api.ServerService
import hata.domain.usecases.LoadRemoteConfigurationUseCase
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers
import javax.inject.Singleton


@Module
@InstallIn(SingletonComponent::class)
object AppModule {

    @Provides
    @Singleton
    fun remoteConfigurationService(): RemoteConfigurationApi {
        return RemoteConfigurationServiceImpl(
            baseUrl = "https://gist.githubusercontent.com/dector/abd9d949a769dc7b0392e7f6529ad148/raw/",
        )
    }

    @Provides
    fun loadRemoteConfigurationUseCase(
        configurationApi: RemoteConfigurationApi,
        @IoDispatcher
        dispatcher: CoroutineDispatcher,
    ): LoadRemoteConfigurationUseCase {
        return LoadRemoteConfigurationUseCase(
            configurationApi = configurationApi,
            dispatcher = dispatcher,
        )
    }

    @Provides
    @Singleton
    fun serverService(): ServerService {
        return RealServerServiceImpl()
    }
}

@Module
@InstallIn(SingletonComponent::class)
object CoroutinesModule {

    @IoDispatcher
    @Provides
    fun ioDispatcher(): CoroutineDispatcher = Dispatchers.IO

    @ComputeDispatcher
    @Provides
    fun computeDispatcher(): CoroutineDispatcher = Dispatchers.Default
}
