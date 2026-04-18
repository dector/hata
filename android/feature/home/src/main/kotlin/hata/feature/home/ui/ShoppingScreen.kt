package hata.feature.home.ui

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.tooling.preview.Preview
import androidx.lifecycle.viewmodel.compose.viewModel
import hata.ui.theme.HataColors
import hata.ui.utils.preview

@Composable
fun ShoppingScreen(
    vm: ShoppingViewModel = viewModel(),
) {
    ShoppingScreenUI(
        onBack = { vm.onDispatch(ShoppingUiAction.Back) },
    )
}

@Composable
private fun ShoppingScreenUI(
    onBack: () -> Unit = {},
) {
    Scaffold(
        containerColor = HataColors.background,
        topBar = {
            ShoppingTopBar(onBack = onBack)
        },
    ) { paddingValues ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues),
            contentAlignment = Alignment.Center,
        ) {
            Text(
                text = "Nothing to buy",
                color = Color.White,
            )
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun ShoppingTopBar(onBack: () -> Unit) {
    TopAppBar(
        title = { Text(text = "Shopping") },
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(
                    imageVector = Icons.AutoMirrored.Filled.ArrowBack,
                    contentDescription = "Back",
                )
            }
        },
        colors = TopAppBarDefaults.topAppBarColors(
            containerColor = HataColors.background,
            titleContentColor = Color.White,
            navigationIconContentColor = Color.White,
        ),
    )
}

@Preview
@Composable
private fun Preview_ShoppingScreen() = preview {
    ShoppingScreenUI()
}
