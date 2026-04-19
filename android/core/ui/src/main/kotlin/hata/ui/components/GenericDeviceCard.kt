package hata.ui.components

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
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AcUnit
import androidx.compose.material.icons.filled.ArrowUpward
import androidx.compose.material.icons.filled.Lightbulb
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.Surface
import androidx.compose.material3.Switch
import androidx.compose.material3.SwitchDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hata.ui.theme.HataColors
import hata.ui.utils.preview


sealed interface DeviceCard {

    data class Generic(
        val id: String,
        val title: String,
        val status: String,
        val isOn: Boolean,
        val canControl: Boolean = true,
        val isAwaitingConfirmation: Boolean = false,
        val icon: Icon = Icon.Default,
    ) : DeviceCard

    enum class Icon {
        Default,
        AC,
        Light,
    }
}

@Composable
fun GenericDeviceCard(
    data: DeviceCard.Generic,
    onToggle: (Boolean) -> Unit = {},
) {
    Surface(
        shape = RoundedCornerShape(24.dp),
        color = if (data.isOn) HataColors.primary else HataColors.surface,
        modifier = Modifier
            .fillMaxWidth()
            .alpha(if (data.canControl) 1f else 0.5f)
            .height(192.dp),
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
                        .background(if (data.isOn) HataColors.primaryContainer else HataColors.surfaceVariant),
                    contentAlignment = Alignment.Center,
                ) {
                    if (data.isAwaitingConfirmation) {
                        CircularProgressIndicator(
                            modifier = Modifier.size(22.dp),
                            color = if (data.isOn) Color.White else HataColors.onSurfaceVariant,
                            strokeWidth = 2.dp,
                        )
                    } else {
                        Icon(
                            imageVector = iconFor(data),
                            contentDescription = null,
                            tint = if (data.isOn) Color.White else HataColors.onSurfaceVariant,
                        )
                    }
                }

                Switch(
                    checked = data.isOn,
                    onCheckedChange = onToggle,
                    enabled = data.canControl && !data.isAwaitingConfirmation,
                    colors = SwitchDefaults.colors(
                        checkedThumbColor = Color.White,
                        checkedTrackColor = HataColors.primaryContainerVariant,
                        uncheckedThumbColor = HataColors.onSurfaceVariant,
                        uncheckedTrackColor = HataColors.surfaceVariant,
                        uncheckedBorderColor = Color.White.copy(alpha = 0.2f),
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

private fun iconFor(data: DeviceCard.Generic): ImageVector {
    return when (data.icon) {
        DeviceCard.Icon.AC -> Icons.Default.AcUnit
        DeviceCard.Icon.Light -> Icons.Default.Lightbulb
        else -> Icons.Default.ArrowUpward
    }
}

@Preview
@Composable
private fun Preview_GenericDeviceCard() = preview(
    pad = 16.dp,
) {
    Column(
        verticalArrangement = Arrangement.spacedBy(16.dp),
    ) {
        GenericDeviceCard(
            DeviceCard.Generic(
                id = "1",
                title = "Air Cooler",
                status = "On",
                isOn = true,
                icon = DeviceCard.Icon.AC,
            ),
        )
        GenericDeviceCard(
            DeviceCard.Generic(
                id = "2",
                title = "Air Cooler",
                status = "Off",
                isOn = false,
                icon = DeviceCard.Icon.AC,
            ),
        )
        GenericDeviceCard(
            DeviceCard.Generic(
                id = "2",
                title = "Air Cooler",
                status = "Off",
                isOn = false,
                icon = DeviceCard.Icon.AC,
                canControl = false,
            ),
        )
    }
}
