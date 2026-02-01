package hata.di

import android.content.Context
import android.content.SharedPreferences
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import hata.core.annotations.NotificationsPreferences
import hata.data.repositories.NotificationsRepository
import hata.data.repositories.NotificationsRepositoryImpl
import javax.inject.Singleton


@Module
@InstallIn(SingletonComponent::class)
object NotificationsModule {

    @Provides
    @Singleton
    @NotificationsPreferences
    fun notificationsPreferences(
        @ApplicationContext context: Context
    ): SharedPreferences {
        return context.getSharedPreferences("hata_notifications", Context.MODE_PRIVATE)
    }

    @Provides
    @Singleton
    fun notificationsRepository(
        @NotificationsPreferences sharedPreferences: SharedPreferences
    ): NotificationsRepository {
        return NotificationsRepositoryImpl(sharedPreferences)
    }
}
