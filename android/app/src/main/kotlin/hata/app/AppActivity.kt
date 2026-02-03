package hata.app

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.runtime.Composable
import androidx.compose.ui.platform.LocalContext
import androidx.core.splashscreen.SplashScreen.Companion.installSplashScreen
import dagger.hilt.android.AndroidEntryPoint
import dagger.hilt.android.EntryPointAccessors
import hata.di.AppRouterEntryPoint
import hata.ui.theme.HataTheme


@AndroidEntryPoint
class AppActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        installSplashScreen()
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            AppContainer()
        }
    }
}

@Composable
fun AppContainer() {
    val entryPoint = EntryPointAccessors.fromApplication(
        LocalContext.current.applicationContext,
        AppRouterEntryPoint::class.java,
    )

    HataTheme {
        AppRouter(
            sessionRepository = entryPoint.sessionRepository(),
            navigator = entryPoint.appNavigator(),
        )
    }
}
