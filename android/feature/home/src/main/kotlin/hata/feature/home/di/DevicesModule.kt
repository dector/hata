package hata.feature.home.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import hata.feature.home.repository.DeviceCacheCleaner
import hata.feature.home.repository.DeviceRepository
import hata.feature.home.repository.DeviceRepositoryImpl
import hata.feature.home.repository.local.UserDevicesDatabaseProvider
import javax.inject.Singleton


@Module
@InstallIn(SingletonComponent::class)
abstract class DevicesModule {

    @Binds
    @Singleton
    abstract fun deviceRepository(
        impl: DeviceRepositoryImpl,
    ): DeviceRepository

    @Binds
    @Singleton
    abstract fun deviceCacheCleaner(
        impl: UserDevicesDatabaseProvider,
    ): DeviceCacheCleaner
}
