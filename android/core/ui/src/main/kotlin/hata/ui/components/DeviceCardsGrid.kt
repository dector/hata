package hata.ui.components

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.GridItemSpan
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
    header: (@Composable () -> Unit)? = null,
    onToggleDevice: (String, Boolean) -> Unit = { _, _ -> },
) {
    LazyVerticalGrid(
        modifier = modifier,
        columns = GridCells.Fixed(columns),
        horizontalArrangement = Arrangement.spacedBy(16.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp),
    ) {
        if (header != null) {
            item(span = { GridItemSpan(maxLineSpan) }) {
                header()
            }
        }

        items(devices) { device ->
            when (device) {
                is DeviceCard.Generic -> GenericDeviceCard(
                    data = device,
                    onToggle = { checked -> onToggleDevice(device.id, checked) },
                )
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
                id = "1",
                title = "Air Cooler",
                status = "On",
                isOn = true,
                icon = DeviceCard.Icon.AC,
            ),
            DeviceCard.Generic(
                id = "2",
                title = "Living Room",
                status = "3 lights on",
                isOn = true,
                icon = DeviceCard.Icon.Light,
            ),
            DeviceCard.Generic(
                id = "3",
                title = "Bedroom",
                status = "Off",
                isOn = false,
                icon = DeviceCard.Icon.Light,
            ),
            DeviceCard.Generic(
                id = "4",
                title = "Kitchen",
                status = "2 lights on",
                isOn = true,
                icon = DeviceCard.Icon.Default,
            ),
        ),
    )
}
