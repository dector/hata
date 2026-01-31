package hata

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
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
    val context = LocalContext.current
    val appContext = context.applicationContext
    val entryPoint = EntryPointAccessors.fromApplication(
        appContext,
        AppRouterEntryPoint::class.java,
    )

    HataTheme {
        Scaffold(
            modifier = Modifier
                .fillMaxSize(),
        ) { innerPadding ->
            AppRouter(
                innerPadding = innerPadding,
                sessionRepository = entryPoint.sessionRepository(),
            )
        }
    }
}
