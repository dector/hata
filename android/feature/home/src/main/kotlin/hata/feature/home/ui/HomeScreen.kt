package hata.feature.home.ui

import androidx.compose.animation.core.LinearEasing
import androidx.compose.animation.core.RepeatMode
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.pulltorefresh.PullToRefreshBox
import androidx.compose.material3.pulltorefresh.rememberPullToRefreshState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.DisposableEffect
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontFamily
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
import hata.ui.components.ConnectionStatus
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
                vm.onDispatch(HomeUiAction.UpdateDevicesStatus())
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
    val isRefreshing = when (state) {
        is HomeUiState.WithData -> state.isSyncing
        HomeUiState.Loading -> true
        HomeUiState.Init, is HomeUiState.Error -> false
    }
    val pullToRefreshState = rememberPullToRefreshState()

    Scaffold { paddingValues ->
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
                connectionStatus = state.toTopBarConnectionStatus(),
                isSyncing = (state as? HomeUiState.WithData)?.isSyncing ?: false,
                hasNotifications = hasUnreadNotifications,
                onAvatarClick = { dispatch(HomeUiAction.NavigateToProfile) },
                onNotificationsClick = { dispatch(HomeUiAction.NavigateToNotifications) },
            )

            PullToRefreshBox(
                modifier = Modifier
                    .fillMaxWidth()
                    .weight(1f),
                state = pullToRefreshState,
                isRefreshing = isRefreshing,
                onRefresh = { dispatch(HomeUiAction.UpdateDevicesStatus(showSyncFeedback = true)) },
                indicator = {},
            ) {
                val revealFraction = pullToRefreshState.distanceFraction.coerceIn(0f, 1f)
                val revealHeight = if (isRefreshing) 32.dp else 32.dp * revealFraction
                val revealTextAlpha = if (isRefreshing) 1f else revealFraction
                val isReadyToRefresh = pullToRefreshState.distanceFraction >= 1f
                val pullHintText = when {
                    isRefreshing -> "Syncing"
                    isReadyToRefresh -> "Release to sync"
                    else -> "Pull to sync"
                }

                Column(
                    modifier = Modifier.fillMaxSize(),
                ) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(revealHeight),
                        contentAlignment = Alignment.Center,
                    ) {
                        if (isRefreshing) {
                            SyncingTextAnimatedDots(
                                color = Color.White,
                                alpha = revealTextAlpha,
                            )
                        } else {
                            Text(
                                text = pullHintText,
                                color = Color.White.copy(alpha = revealTextAlpha),
                                fontSize = 14.sp,
                                fontWeight = FontWeight.Medium,
                            )
                        }
                    }

                    StateSection(
                        state = state,
                        header = {
                            Spacer(modifier = Modifier.height(16.dp))
                            WelcomeSection()
                            Spacer(modifier = Modifier.height(32.dp))
                        },
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
    }
}

@Composable
private fun SyncingTextAnimatedDots(
    modifier: Modifier = Modifier,
    color: Color = Color.White,
    alpha: Float = 1f,
) {
    val transition = rememberInfiniteTransition(label = "syncing_text_dots")
    val dotsProgress by transition.animateFloat(
        initialValue = 0f,
        targetValue = 3f,
        animationSpec = infiniteRepeatable(
            animation = tween(durationMillis = 900, easing = LinearEasing),
            repeatMode = RepeatMode.Restart,
        ),
        label = "syncing_text_dots_count",
    )
    val dotsCount = (dotsProgress.toInt() % 3) + 1
    val dotsText = ".".repeat(dotsCount).padEnd(3, ' ')

    Text(
        modifier = modifier,
        text = "Syncing $dotsText",
        color = color.copy(alpha = alpha),
        fontSize = 14.sp,
        fontWeight = FontWeight.Medium,
        fontFamily = FontFamily.Monospace,
    )
}

@Composable
private fun StateSection(
    state: HomeUiState,
    header: @Composable () -> Unit = {},
    onToggleDevice: (String, Boolean) -> Unit = { _, _ -> },
) {
    when (state) {
        is HomeUiState.Init -> {}

        is HomeUiState.Loading -> {
            Column(
                modifier = Modifier.fillMaxSize(),
            ) {
                header()
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center,
                ) {
                    CircularProgressIndicator(color = Color.White)
                }

            }
        }

        is HomeUiState.WithData -> {
            DeviceCardsGrid(
                modifier = Modifier.padding(horizontal = 16.dp),
                devices = state.data.devices,
                header = header,
                onToggleDevice = onToggleDevice,
            )
        }

        is HomeUiState.Error -> {
            Column(
                modifier = Modifier.fillMaxSize(),
            ) {
                header()
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
                connectionStatus = ServerConnectionStatus.Online,
            ),
        ),
    )
}

private fun HomeUiState.toTopBarConnectionStatus(): ConnectionStatus =
    when (this) {
        is HomeUiState.WithData -> data.connectionStatus.toTopBarConnectionStatus()
        is HomeUiState.Error -> connectionStatus.toTopBarConnectionStatus()
        HomeUiState.Loading -> ServerConnectionStatus.Unknown.toTopBarConnectionStatus()
        HomeUiState.Init -> ServerConnectionStatus.Unknown.toTopBarConnectionStatus()
    }

private fun ServerConnectionStatus.toTopBarConnectionStatus(): ConnectionStatus =
    when (this) {
        ServerConnectionStatus.Unknown -> ConnectionStatus.Unknown
        ServerConnectionStatus.Online -> ConnectionStatus.Online
        ServerConnectionStatus.Offline -> ConnectionStatus.Offline
    }
