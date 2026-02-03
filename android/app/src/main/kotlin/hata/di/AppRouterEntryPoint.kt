package hata.di

import dagger.hilt.EntryPoint
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import hata.feature.session.repository.SessionRepository
import hata.navigation.internal.AppNavigator


@EntryPoint
@InstallIn(SingletonComponent::class)
interface AppRouterEntryPoint {
    fun sessionRepository(): SessionRepository
    fun appNavigator(): AppNavigator
}
