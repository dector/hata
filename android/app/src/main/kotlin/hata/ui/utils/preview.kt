package hata

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import hata.ui.theme.HataTheme

@Composable
fun preview(
    pad: Dp = 0.dp,
    content: @Composable () -> Unit,
) {
    HataTheme {
        Box(
            modifier = Modifier.Companion
                .background(MaterialTheme.colorScheme.background)
                .padding(pad),
        ) {
            content()
        }
    }
}
