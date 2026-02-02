package hata.ui.components

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import hata.ui.utils.preview


@Composable
fun DeviceCardsGrid(
    modifier: Modifier = Modifier,
    devices: List<DeviceCard>,
    columns: Int = 2,
) {
    LazyVerticalGrid(
        modifier = modifier,
        columns = GridCells.Fixed(columns),
        horizontalArrangement = Arrangement.spacedBy(16.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp),
    ) {
        items(devices) { device ->
            when (device) {
                is DeviceCard.Generic -> GenericDeviceCard(device)
            }
        }
    }
}

@Preview
@Composable
private fun Preview_DeviceCardsGrid() = preview(
    pad = 16.dp,
) {
    DeviceCardsGrid(
        devices = listOf(
            DeviceCard.Generic(
                title = "Air Cooler",
                status = "On",
                isOn = true,
                icon = DeviceCard.Icon.AC,
            ),
            DeviceCard.Generic(
                title = "Living Room",
                status = "3 lights on",
                isOn = true,
                icon = DeviceCard.Icon.Light,
            ),
            DeviceCard.Generic(
                title = "Bedroom",
                status = "Off",
                isOn = false,
                icon = DeviceCard.Icon.Light,
            ),
            DeviceCard.Generic(
                title = "Kitchen",
                status = "2 lights on",
                isOn = true,
                icon = DeviceCard.Icon.Default,
            ),
        ),
    )
}
