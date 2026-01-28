package hata.ui.home

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hata.ui.components.DeviceCard
import hata.ui.components.DeviceCardsGrid
import hata.ui.components.TopBar
import hata.ui.theme.Colors
import hata.ui.utils.preview


@Composable
fun HomeScreen() {
    HomeScreenUI()
}

@Composable
@OptIn(ExperimentalMaterial3Api::class)
private fun HomeScreenUI() {
    Scaffold(
//        bottomBar = { BottomNavigationBar() },
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .background(Colors.mainBg),
        ) {
            TopBar(
                name = "Home",
            )

            Spacer(modifier = Modifier.height(16.dp))

            // Welcome text
            WelcomeSection()

            Spacer(modifier = Modifier.height(32.dp))

            val devices = listOf<DeviceCard>(
                DeviceCard.Generic(
                    title = "Air Cooler",
                    status = "On",
                    isOn = true,
                    icon = DeviceCard.Icon.AC,
                ),
                DeviceCard.Generic(
                    title = "Lamp",
                    status = "Off",
                    isOn = false,
                    icon = DeviceCard.Icon.Light,
                ),
            )
            DeviceCardsGrid(
                modifier = Modifier.padding(horizontal = 16.dp),
                devices = devices,
            )
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
    HomeScreenUI()
}
