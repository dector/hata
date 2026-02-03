package hata.feature.session.di

import android.content.Context
import android.content.SharedPreferences
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import hata.core.annotations.AuthPreferences
import hata.feature.session.repository.SessionRepository
import hata.feature.session.repository.SessionRepositoryImpl
import javax.inject.Singleton


@Module
@InstallIn(SingletonComponent::class)
object SessionModule {

    @Provides
    @Singleton
    @AuthPreferences
    fun authPreferences(
        @ApplicationContext context: Context,
    ): SharedPreferences {
        return context.getSharedPreferences("hata_auth", Context.MODE_PRIVATE)
    }

    @Provides
    @Singleton
    fun sessionRepository(
        @AuthPreferences sharedPreferences: SharedPreferences,
    ): SessionRepository {
        return SessionRepositoryImpl(sharedPreferences)
    }
}
