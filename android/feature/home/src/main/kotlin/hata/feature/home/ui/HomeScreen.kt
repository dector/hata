package hata.feature.home.ui

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.DisposableEffect
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.LifecycleEventObserver
import androidx.lifecycle.compose.LocalLifecycleOwner
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import hata.feature.home.domain.NewState
import hata.ui.components.DeviceCard
import hata.ui.components.DeviceCardsGrid
import hata.ui.components.TopBar
import hata.ui.theme.HataColors
import hata.ui.utils.preview


@Composable
fun HomeScreen(
    vm: HomeViewModel = viewModel(),
) {
    val state by vm.uiState.collectAsStateWithLifecycle()
    val hasUnreadNotifications by vm.hasUnreadNotifications.collectAsStateWithLifecycle()
    val lifecycleOwner = LocalLifecycleOwner.current

    LaunchedEffect(Unit) {
        vm.onDispatch(HomeUiAction.Init)
    }

    DisposableEffect(lifecycleOwner) {
        val observer = LifecycleEventObserver { _, event ->
            if (event == Lifecycle.Event.ON_START) {
                vm.onDispatch(HomeUiAction.UpdateDevicesStatus)
            }
        }

        lifecycleOwner.lifecycle.addObserver(observer)
        onDispose { lifecycleOwner.lifecycle.removeObserver(observer) }
    }

    HomeScreenUI(
        state = state,
        hasUnreadNotifications = hasUnreadNotifications,
        dispatch = vm::onDispatch,
    )
}

@Composable
@OptIn(ExperimentalMaterial3Api::class)
private fun HomeScreenUI(
    state: HomeUiState,
    hasUnreadNotifications: Boolean = false,
    dispatch: (HomeUiAction) -> Unit = {},
) {
    Scaffold(
//        bottomBar = { BottomNavigationBar() },
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .background(HataColors.background)
                .padding(paddingValues),
        ) {
            val homeName = when (state) {
                is HomeUiState.WithData -> state.data.home.name
                else -> "Home"
            }
            TopBar(
                name = homeName,
                hasNotifications = hasUnreadNotifications,
                onAvatarClick = { dispatch(HomeUiAction.NavigateToProfile) },
                onNotificationsClick = { dispatch(HomeUiAction.NavigateToNotifications) },
            )

            Spacer(modifier = Modifier.height(16.dp))

            WelcomeSection()

            Spacer(modifier = Modifier.height(32.dp))

            StateSection(
                state = state,
                onToggleDevice = { deviceId, isEnabled ->
                    dispatch(
                        HomeUiAction.ToggleDevice(
                            deviceId = deviceId,
                            newState = if (isEnabled) NewState.On else NewState.Off,
                        ),
                    )
                },
            )
        }
    }
}

@Composable
private fun StateSection(
    state: HomeUiState,
    onToggleDevice: (String, Boolean) -> Unit = { _, _ -> },
) {
    when (state) {
        is HomeUiState.Init -> {}

        is HomeUiState.Loading -> {
            Box(
                modifier = Modifier.fillMaxSize(),
                contentAlignment = Alignment.Center,
            ) {
                CircularProgressIndicator(color = Color.White)

            }
        }

        is HomeUiState.WithData -> {
            DeviceCardsGrid(
                modifier = Modifier.padding(horizontal = 16.dp),
                devices = state.data.devices,
                onToggleDevice = onToggleDevice,
            )
        }

        is HomeUiState.Error -> {
            Box(
                modifier = Modifier.fillMaxSize(),
                contentAlignment = Alignment.Center,
            ) {
                Text(
                    text = "Error: ${state.message}",
                    color = Color.Red,
                    modifier = Modifier.padding(16.dp),
                )
            }
        }
    }
}

@Composable
private fun WelcomeSection() {
    Text(
        modifier = Modifier.padding(horizontal = 16.dp),
        text = "Welcome home,\nDan",
        color = Color.White,
        fontSize = 36.sp,
        fontWeight = FontWeight.Bold,
        lineHeight = 42.sp,
    )
}

/*
@Composable
private fun BottomNavigationBar() {
    NavigationBar(
        containerColor = Color(0xFF1E2423),
        contentColor = Color.White,
    ) {
        NavigationBarItem(
            icon = {
                Icon(
                    imageVector = Icons.Default.Home,
                    contentDescription = "Home",
                )
            },
            label = { Text("Home") },
            selected = true,
            onClick = {},
            colors = NavigationBarItemDefaults.colors(
                selectedIconColor = Color(0xFF8FB899),
                selectedTextColor = Color(0xFF8FB899),
                indicatorColor = Color(0xFF2D3533),
                unselectedIconColor = Color(0xFF6B7875),
                unselectedTextColor = Color(0xFF6B7875),
            ),
        )
        NavigationBarItem(
            icon = {
                Icon(
                    imageVector = Icons.Default.Dashboard,
                    contentDescription = "Scenes",
                )
            },
            label = { Text("Scenes") },
            selected = false,
            onClick = {},
            colors = NavigationBarItemDefaults.colors(
                selectedIconColor = Color(0xFF8FB899),
                selectedTextColor = Color(0xFF8FB899),
                indicatorColor = Color(0xFF2D3533),
                unselectedIconColor = Color(0xFF6B7875),
                unselectedTextColor = Color(0xFF6B7875),
            ),
        )
        NavigationBarItem(
            icon = {
                Icon(
                    imageVector = Icons.Default.Schedule,
                    contentDescription = "Auto",
                )
            },
            label = { Text("Auto") },
            selected = false,
            onClick = {},
            colors = NavigationBarItemDefaults.colors(
                selectedIconColor = Color(0xFF8FB899),
                selectedTextColor = Color(0xFF8FB899),
                indicatorColor = Color(0xFF2D3533),
                unselectedIconColor = Color(0xFF6B7875),
                unselectedTextColor = Color(0xFF6B7875),
            ),
        )
        NavigationBarItem(
            icon = {
                Icon(
                    imageVector = Icons.Default.Settings,
                    contentDescription = "Settings",
                )
            },
            label = { Text("Settings") },
            selected = false,
            onClick = {},
            colors = NavigationBarItemDefaults.colors(
                selectedIconColor = Color(0xFF8FB899),
                selectedTextColor = Color(0xFF8FB899),
                indicatorColor = Color(0xFF2D3533),
                unselectedIconColor = Color(0xFF6B7875),
                unselectedTextColor = Color(0xFF6B7875),
            ),
        )
    }
}
*/

@Preview
@Composable
private fun Preview_HomeScreen() = preview {
    HomeScreenUI(
        state = HomeUiState.WithData(
            data = HomeDisplayData(
                home = HomeDisplay(name = "My Home"),
                devices = listOf(
                    DeviceCard.Generic(
                        id = "1",
                        title = "Living Room",
                        status = "3 lights on",
                        isOn = true,
                        icon = DeviceCard.Icon.Light,
                    ),
                    DeviceCard.Generic(
                        id = "2",
                        title = "Bedroom",
                        status = "Off",
                        isOn = false,
                        icon = DeviceCard.Icon.Light,
                    ),
                    DeviceCard.Generic(
                        id = "3",
                        title = "Air Cooler",
                        status = "On",
                        isOn = true,
                        icon = DeviceCard.Icon.AC,
                    ),
                    DeviceCard.Generic(
                        id = "4",
                        title = "Kitchen",
                        status = "Ready",
                        isOn = false,
                        icon = DeviceCard.Icon.Default,
                    ),
                ),
            ),
        ),
    )
}
