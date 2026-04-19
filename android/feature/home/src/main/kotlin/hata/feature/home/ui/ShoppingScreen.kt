package hata.feature.home.ui

import androidx.compose.animation.core.LinearEasing
import androidx.compose.animation.core.RepeatMode
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Done
import androidx.compose.material.icons.filled.RadioButtonUnchecked
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SwipeToDismissBox
import androidx.compose.material3.pulltorefresh.PullToRefreshBox
import androidx.compose.material3.pulltorefresh.rememberPullToRefreshState
import androidx.compose.material3.SwipeToDismissBoxValue
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.material3.rememberSwipeToDismissBoxState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import hata.feature.home.domain.ShoppingItem
import hata.ui.theme.HataColors
import hata.ui.utils.preview

@Composable
fun ShoppingScreen(
    vm: ShoppingViewModel = viewModel(),
) {
    val state by vm.uiState.collectAsStateWithLifecycle()

    LaunchedEffect(Unit) {
        vm.onDispatch(ShoppingUiAction.Init)
    }

    ShoppingScreenUI(
        state = state,
        dispatch = vm::onDispatch,
    )
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun ShoppingScreenUI(
    state: ShoppingUiState,
    dispatch: (ShoppingUiAction) -> Unit = {},
) {
    val isRefreshing = when (state) {
        ShoppingUiState.Loading -> true
        is ShoppingUiState.Loaded -> state.isRefreshing
        ShoppingUiState.Init, is ShoppingUiState.Error -> false
    }
    val pullToRefreshState = rememberPullToRefreshState()

    Scaffold(
        containerColor = HataColors.background,
        topBar = {
            ShoppingTopBar(onBack = { dispatch(ShoppingUiAction.Back) })
        },
    ) { paddingValues ->
        PullToRefreshBox(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues),
            state = pullToRefreshState,
            isRefreshing = isRefreshing,
            onRefresh = { dispatch(ShoppingUiAction.Refresh) },
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

                Box(
                    modifier = Modifier.fillMaxSize(),
                ) {
                    when (state) {
                        ShoppingUiState.Init,
                        ShoppingUiState.Loading,
                        -> ShoppingLoadingState()

                        is ShoppingUiState.Loaded -> {
                            if (state.items.isEmpty()) {
                                ShoppingEmptyState()
                            } else {
                                ShoppingItemsList(
                                    items = state.items,
                                    errorMessage = state.errorMessage,
                                    onMarkPurchased = { itemId ->
                                        dispatch(ShoppingUiAction.MarkPurchased(itemId))
                                    },
                                )
                            }
                        }

                        is ShoppingUiState.Error -> {
                            ShoppingErrorState(
                                message = state.message,
                                onRetry = { dispatch(ShoppingUiAction.Retry) },
                            )
                        }
                    }
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
private fun ShoppingLoadingState() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center,
    ) {
        CircularProgressIndicator(color = HataColors.primary)
    }
}

@Composable
private fun ShoppingEmptyState() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center,
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            Text(
                text = "Nothing to buy",
                color = Color.White,
                fontSize = 18.sp,
                fontWeight = FontWeight.SemiBold,
            )
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = "Your default shopping list is empty",
                color = HataColors.onSurfaceVariant,
            )
        }
    }
}

@Composable
private fun ShoppingErrorState(
    message: String,
    onRetry: () -> Unit,
) {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center,
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            modifier = Modifier.padding(horizontal = 24.dp),
        ) {
            Text(
                text = message,
                color = HataColors.error,
            )
            Spacer(modifier = Modifier.height(16.dp))
            Button(onClick = onRetry) {
                Text(text = "Retry")
            }
        }
    }
}

@Composable
private fun ShoppingItemsList(
    items: List<ShoppingItem>,
    errorMessage: String?,
    onMarkPurchased: (String) -> Unit,
) {
    LazyColumn(
        modifier = Modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(8.dp),
    ) {
        if (!errorMessage.isNullOrBlank()) {
            item {
                Text(
                    text = errorMessage,
                    color = HataColors.error,
                    modifier = Modifier.padding(horizontal = 16.dp),
                )
            }
        }

        items(items, key = { it.id }) { item ->
            ShoppingItemRow(
                item = item,
                onMarkPurchased = { onMarkPurchased(item.id) },
            )
        }
        item { Spacer(modifier = Modifier.height(16.dp)) }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun ShoppingItemRow(
    item: ShoppingItem,
    onMarkPurchased: () -> Unit,
) {
    val dismissState = rememberSwipeToDismissBoxState(
        confirmValueChange = { value ->
            if (value == SwipeToDismissBoxValue.StartToEnd && !item.isChecked) {
                onMarkPurchased()
            }
            false
        },
    )

    SwipeToDismissBox(
        state = dismissState,
        enableDismissFromStartToEnd = !item.isChecked,
        enableDismissFromEndToStart = false,
        backgroundContent = {
            Row(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(horizontal = 32.dp),
                verticalAlignment = Alignment.CenterVertically,
            ) {
                Icon(
                    imageVector = Icons.Default.Done,
                    contentDescription = null,
                    tint = HataColors.primary,
                )
            }
        },
    ) {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 16.dp),
            shape = RoundedCornerShape(12.dp),
            colors = CardDefaults.cardColors(
                containerColor = HataColors.surface,
            ),
        ) {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 16.dp, vertical = 14.dp),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(12.dp),
            ) {
                Icon(
                    imageVector = if (item.isChecked) {
                        Icons.Default.CheckCircle
                    } else {
                        Icons.Default.RadioButtonUnchecked
                    },
                    contentDescription = null,
                    tint = if (item.isChecked) HataColors.primary else HataColors.onSurfaceVariant,
                )
                Text(
                    text = item.name,
                    color = if (item.isChecked) HataColors.onSurfaceVariant else Color.White,
                    textDecoration = if (item.isChecked) TextDecoration.LineThrough else TextDecoration.None,
                )
            }
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
private fun Preview_ShoppingScreen_Loaded() = preview {
    ShoppingScreenUI(
        state = ShoppingUiState.Loaded(
            items = listOf(
                ShoppingItem(id = "1", name = "Milk", isChecked = false),
                ShoppingItem(id = "2", name = "Eggs", isChecked = true),
            ),
        ),
    )
}

@Preview
@Composable
private fun Preview_ShoppingScreen_Empty() = preview {
    ShoppingScreenUI(
        state = ShoppingUiState.Loaded(items = emptyList()),
    )
}

@Preview
@Composable
private fun Preview_ShoppingScreen_Error() = preview {
    ShoppingScreenUI(
        state = ShoppingUiState.Error(message = "Failed to load shopping list"),
    )
}
