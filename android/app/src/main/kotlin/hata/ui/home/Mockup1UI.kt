package hata.ui.home

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AcUnit
import androidx.compose.material.icons.filled.Dashboard
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.Lightbulb
import androidx.compose.material.icons.filled.Lock
import androidx.compose.material.icons.filled.MoreVert
import androidx.compose.material.icons.filled.Schedule
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material.icons.filled.Shield
import androidx.compose.material.icons.filled.Thermostat
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Switch
import androidx.compose.material3.SwitchDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hata.ui.components.TopBar
import hata.ui.theme.Colors
import hata.ui.utils.preview


@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun Mockup1UI() {
    Scaffold(
        bottomBar = { BottomNavigationBar() },
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .background(Colors.mainBg)
                .padding(paddingValues),
        ) {
            TopBar(
                name = "Home",
            )

            Spacer(modifier = Modifier.height(16.dp))

            // Welcome text
            WelcomeSection()

            Spacer(modifier = Modifier.height(32.dp))

            // Status chips
            StatusChips()

            Spacer(modifier = Modifier.height(24.dp))

            // Device cards grid
            DeviceCardsGrid()
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

@Composable
private fun StatusChips() {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp),
        horizontalArrangement = Arrangement.spacedBy(12.dp),
    ) {
        StatusChip(
            icon = Icons.Default.Lightbulb,
            text = "3 lights on",
            iconTint = Color(0xFF8FB899),
        )
        StatusChip(
            icon = Icons.Default.Thermostat,
            text = "AC 22°C",
            iconTint = Color(0xFF8FB899),
        )
        StatusChip(
            icon = Icons.Default.Shield,
            text = "Active",
            iconTint = Color(0xFF8FB899),
        )
    }
}

@Composable
private fun StatusChip(
    icon: ImageVector,
    text: String,
    iconTint: Color,
) {
    Surface(
        shape = RoundedCornerShape(24.dp),
        color = Color(0xFF2D3533),
        modifier = Modifier.height(48.dp),
    ) {
        Row(
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 12.dp),
            horizontalArrangement = Arrangement.spacedBy(8.dp),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Icon(
                imageVector = icon,
                contentDescription = null,
                tint = iconTint,
                modifier = Modifier.size(20.dp),
            )
            Text(
                text = text,
                color = Color.White,
                fontSize = 14.sp,
            )
        }
    }
}

@Composable
private fun DeviceCardsGrid() {
    val devices = listOf(
        DeviceCardData.LivingRoomLamp(),
        DeviceCardData.KitchenAC(),
        DeviceCardData.FrontDoor(),
        DeviceCardData.AllFloorLights(),
    )

    LazyVerticalGrid(
        modifier = Modifier.padding(horizontal = 16.dp),
        columns = GridCells.Fixed(2),
        horizontalArrangement = Arrangement.spacedBy(16.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp),
    ) {
        items(devices) { device ->
            when (device) {
                is DeviceCardData.LivingRoomLamp -> LampCard(device)
                is DeviceCardData.KitchenAC -> ACCard(device)
                is DeviceCardData.FrontDoor -> DoorCard(device)
                is DeviceCardData.AllFloorLights -> LightsGroupCard(device)
            }
        }
    }
}

@Composable
private fun LampCard(data: DeviceCardData.LivingRoomLamp) {
    Surface(
        shape = RoundedCornerShape(24.dp),
        color = Color(0xFF2D3533),
        modifier = Modifier
            .fillMaxWidth()
            .height(280.dp),
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(20.dp),
        ) {
            Box(
                modifier = Modifier
                    .size(48.dp)
                    .clip(CircleShape)
                    .background(Color(0xFF3D4845)),
                contentAlignment = Alignment.Center,
            ) {
                Icon(
                    imageVector = Icons.Default.Lightbulb,
                    contentDescription = null,
                    tint = Color(0xFF6B7875),
                )
            }

            Spacer(modifier = Modifier.weight(1f))

            Text(
                text = "${data.brightness}%",
                color = Color.White.copy(alpha = 0.9f),
                fontSize = 28.sp,
                fontWeight = FontWeight.Medium,
            )

            Spacer(modifier = Modifier.height(24.dp))

            Text(
                text = data.title,
                color = Color.White,
                fontSize = 20.sp,
                fontWeight = FontWeight.Medium,
            )
            Text(
                text = data.status,
                color = Color(0xFF6B7875),
                fontSize = 14.sp,
            )
        }
    }
}

@Composable
private fun ACCard(data: DeviceCardData.KitchenAC) {
    Surface(
        shape = RoundedCornerShape(24.dp),
        color = Color(0xFF8FB899),
        modifier = Modifier
            .fillMaxWidth()
            .height(280.dp),
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(20.dp),
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.Top,
            ) {
                Box(
                    modifier = Modifier
                        .size(48.dp)
                        .clip(CircleShape)
                        .background(Color(0xFFA8C5B0)),
                    contentAlignment = Alignment.Center,
                ) {
                    Icon(
                        imageVector = Icons.Default.AcUnit,
                        contentDescription = null,
                        tint = Color.White,
                    )
                }

                Switch(
                    checked = data.isOn,
                    onCheckedChange = {},
                    colors = SwitchDefaults.colors(
                        checkedThumbColor = Color.White,
                        checkedTrackColor = Color(0xFFC8DFD0),
                        uncheckedThumbColor = Color(0xFF6B7875),
                        uncheckedTrackColor = Color(0xFF3D4845),
                    ),
                )
            }

            Spacer(modifier = Modifier.weight(1f))

            Text(
                text = data.title,
                color = Color.White,
                fontSize = 20.sp,
                fontWeight = FontWeight.Medium,
            )
            Text(
                text = data.status,
                color = Color.White.copy(alpha = 0.85f),
                fontSize = 14.sp,
            )
        }
    }
}

@Composable
private fun DoorCard(data: DeviceCardData.FrontDoor) {
    Surface(
        shape = RoundedCornerShape(24.dp),
        color = Color(0xFF2D3533),
        modifier = Modifier
            .fillMaxWidth()
            .height(280.dp),
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(20.dp),
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.Top,
            ) {
                Box(
                    modifier = Modifier
                        .size(48.dp)
                        .clip(CircleShape)
                        .background(Color(0xFF3D4845)),
                    contentAlignment = Alignment.Center,
                ) {
                    Icon(
                        imageVector = Icons.Default.Lock,
                        contentDescription = null,
                        tint = Color(0xFF6B7875),
                    )
                }

                Icon(
                    imageVector = Icons.Default.MoreVert,
                    contentDescription = "More options",
                    tint = Color(0xFF6B7875),
                )
            }

            Spacer(modifier = Modifier.weight(1f))

            Text(
                text = data.title,
                color = Color.White,
                fontSize = 20.sp,
                fontWeight = FontWeight.Medium,
            )
            Text(
                text = data.status,
                color = Color(0xFF6B7875),
                fontSize = 14.sp,
            )
        }
    }
}

@Composable
private fun LightsGroupCard(data: DeviceCardData.AllFloorLights) {
    Surface(
        shape = RoundedCornerShape(24.dp),
        color = Color(0xFF2D3533),
        modifier = Modifier
            .fillMaxWidth()
            .height(280.dp),
    ) {
        Box(modifier = Modifier.fillMaxSize()) {
            // Background group text
            Text(
                text = "GROUP",
                color = Color(0xFF3D4845),
                fontSize = 48.sp,
                fontWeight = FontWeight.Bold,
                modifier = Modifier
                    .align(Alignment.CenterEnd)
                    .padding(end = 16.dp),
            )

            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(20.dp),
            ) {
                Box(
                    modifier = Modifier
                        .size(48.dp)
                        .clip(CircleShape)
                        .background(Color(0xFF3D4845)),
                    contentAlignment = Alignment.Center,
                ) {
                    Icon(
                        imageVector = Icons.Default.Lightbulb,
                        contentDescription = null,
                        tint = Color(0xFF6B7875),
                    )
                }

                Spacer(modifier = Modifier.weight(1f))

                Text(
                    text = data.title,
                    color = Color.White,
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Medium,
                )
                Text(
                    text = data.status,
                    color = Color(0xFF6B7875),
                    fontSize = 14.sp,
                )
            }
        }
    }
}

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

// Data models
sealed class DeviceCardData {
    data class LivingRoomLamp(
        val title: String = "Living Room\nLamp",
        val status: String = "On",
        val brightness: Int = 80,
    ) : DeviceCardData()

    data class KitchenAC(
        val title: String = "Kitchen\nAC",
        val status: String = "Cooling to 22°C",
        val isOn: Boolean = true,
    ) : DeviceCardData()

    data class FrontDoor(
        val title: String = "Front\nDoor",
        val status: String = "Locked",
    ) : DeviceCardData()

    data class AllFloorLights(
        val title: String = "All Floor\nLights",
        val status: String = "2/5 On",
    ) : DeviceCardData()
}

@Preview
@Composable
private fun Preview_Mockup1UI() = preview {
    Mockup1UI()
}
