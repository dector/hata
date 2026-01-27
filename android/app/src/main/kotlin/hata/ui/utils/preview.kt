package hata.ui.utils

import android.annotation.SuppressLint
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import hata.ui.theme.HataTheme


@Composable
@SuppressLint("ComposableNaming")
fun preview(
    pad: Dp = 0.dp,
    bg: Color? = null,
    content: @Composable () -> Unit,
) {
    HataTheme {
        Box(
            modifier = Modifier
                .background(bg ?: MaterialTheme.colorScheme.background)
                .padding(pad),
        ) {
            content()
        }
    }
}
