package hata.ui.theme

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

data class ColorInfo(
    val color: Color,
    val name: String,
    val usage: String,
)

@Composable
fun ColorPalettePreview() {
    val semanticColors = listOf(
        ColorInfo(HataColors.primary, "primary", "Main brand color, buttons, active states"),
        ColorInfo(HataColors.onPrimary, "onPrimary", "Text/icons on primary color"),
        ColorInfo(HataColors.background, "background", "Main app background"),
        ColorInfo(HataColors.surface, "surface", "Card backgrounds, elevated surfaces"),
        ColorInfo(HataColors.surfaceVariant, "surfaceVariant", "Icon container backgrounds"),
        ColorInfo(HataColors.onSurface, "onSurface", "Primary text color"),
        ColorInfo(HataColors.onSurfaceVariant, "onSurfaceVariant", "Secondary text, inactive states"),
        ColorInfo(HataColors.primaryContainer, "primaryContainer", "Active state backgrounds, light primary"),
        ColorInfo(HataColors.primaryContainerVariant, "primaryContainerVariant", "Switch track (checked state)"),
        ColorInfo(HataColors.error, "error", "Error messages and text"),
        ColorInfo(HataColors.errorVariant, "errorVariant", "Notification indicators, alerts"),
    )

    val accentColors = listOf(
        ColorInfo(HataAccentColors.avatarBackground, "avatarBackground", "Avatar background color"),
        ColorInfo(HataAccentColors.avatarIcon, "avatarIcon", "Avatar icon tint color"),
    )

    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .background(HataColors.background)
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp),
    ) {
        item {
            SectionHeader("Core Semantic Colors")
        }
        items(semanticColors) { colorInfo ->
            ColorItem(
                color = colorInfo.color,
                name = colorInfo.name,
                usage = colorInfo.usage,
            )
        }

        item {
            SectionHeader("Accent Colors")
        }
        items(accentColors) { colorInfo ->
            ColorItem(
                color = colorInfo.color,
                name = colorInfo.name,
                usage = colorInfo.usage,
            )
        }
    }
}

@Composable
private fun SectionHeader(title: String) {
    Text(
        text = title,
        color = Color.White,
        fontSize = 18.sp,
        fontWeight = FontWeight.Bold,
        modifier = Modifier.padding(vertical = 12.dp),
    )
}

@Composable
private fun ColorItem(
    color: Color,
    name: String,
    usage: String,
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .background(HataColors.surface, RoundedCornerShape(12.dp))
            .padding(16.dp),
        horizontalArrangement = Arrangement.spacedBy(16.dp),
    ) {
        // Color swatch
        Box(
            modifier = Modifier
                .size(48.dp)
                .background(color, RoundedCornerShape(8.dp)),
        )

        Column(
            modifier = Modifier.weight(1f),
            verticalArrangement = Arrangement.spacedBy(4.dp),
        ) {
            Text(
                text = name,
                color = Color.White,
                fontSize = 16.sp,
                fontWeight = FontWeight.Bold,
            )
            Text(
                text = colorToHex(color),
                color = HataColors.onSurfaceVariant,
                fontSize = 14.sp,
                fontWeight = FontWeight.Medium,
            )
            Text(
                text = usage,
                color = HataColors.onSurfaceVariant,
                fontSize = 12.sp,
            )
        }
    }
}

private fun colorToHex(color: Color): String {
    val red = (color.red * 255).toInt()
    val green = (color.green * 255).toInt()
    val blue = (color.blue * 255).toInt()
    return "#%02X%02X%02X".format(red, green, blue)
}

@Preview(showBackground = true)
@Composable
private fun Preview_ColorPalette() {
    HataTheme {
        ColorPalettePreview()
    }
}
