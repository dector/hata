package hata.navigation.di

import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import hata.navigation.Navigator
import hata.navigation.internal.AppNavigator
import hata.navigation.internal.GlobalNavigator
import javax.inject.Singleton


@Module
@InstallIn(SingletonComponent::class)
object NavigationModule {

    @Provides
    @Singleton
    fun navigator(
        instance: GlobalNavigator,
    ): Navigator = instance.exposed

    @Provides
    @Singleton
    fun appNavigator(
        instance: GlobalNavigator,
    ): AppNavigator = instance

    @Provides
    @Singleton
    fun globalNavigator() = GlobalNavigator()
}
