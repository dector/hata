package space.dector.hata.watch

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.pager.HorizontalPager
import androidx.compose.foundation.pager.rememberPagerState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.ExperimentalMaterialApi
import androidx.compose.material.pullrefresh.PullRefreshIndicator
import androidx.compose.material.pullrefresh.pullRefresh
import androidx.compose.material.pullrefresh.rememberPullRefreshState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.wear.compose.material.Button
import androidx.wear.compose.material.CircularProgressIndicator
import androidx.wear.compose.material.MaterialTheme
import androidx.wear.compose.material.PositionIndicator
import androidx.wear.compose.material.ScalingLazyColumn
import androidx.wear.compose.material.Text
import androidx.wear.compose.material.rememberScalingLazyListState
import com.google.android.gms.wearable.CapabilityClient
import com.google.android.gms.wearable.MessageClient
import com.google.android.gms.wearable.MessageEvent
import com.google.android.gms.wearable.Wearable
import org.json.JSONObject
import java.nio.charset.StandardCharsets
import java.time.Duration
import java.time.Instant

private const val PhoneCapability = "hata_phone_app"
private const val WeatherRequestPath = "/watch/weather/request"
private const val WeatherUpdatePath = "/phone/weather/update"
private const val WeatherErrorPath = "/phone/weather/error"
private const val DevicesRequestPath = "/watch/devices/request"
private const val DeviceTogglePath = "/watch/device/toggle"
private const val DevicesUpdatePath = "/phone/devices/update"
private const val DevicesErrorPath = "/phone/devices/error"

class MainActivity : ComponentActivity(), MessageClient.OnMessageReceivedListener {
    private var weatherState by mutableStateOf<WeatherScreenState>(WeatherScreenState.Connecting)
    private var devicesState by mutableStateOf<DevicesScreenState>(DevicesScreenState.Loading)
    private var togglingDeviceIds by mutableStateOf<Set<String>>(emptySet())

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        setContent {
            HataWatchApp(
                weatherState = weatherState,
                devicesState = devicesState,
                togglingDeviceIds = togglingDeviceIds,
                onRetryWeather = ::requestWeather,
                onRetryDevices = ::requestDevices,
                onDeviceTap = ::toggleDevice,
            )
        }
    }

    override fun onResume() {
        super.onResume()
        Wearable.getMessageClient(this).addListener(this)
        requestWeather()
        requestDevices()
    }

    override fun onPause() {
        Wearable.getMessageClient(this).removeListener(this)
        super.onPause()
    }

    override fun onMessageReceived(event: MessageEvent) {
        val body = event.data.toString(StandardCharsets.UTF_8)
        runOnUiThread {
            when (event.path) {
                WeatherUpdatePath -> weatherState = parseWeather(body)
                WeatherErrorPath -> weatherState = WeatherScreenState.Error(errorMessage(body, "Weather unavailable"))
                DevicesUpdatePath -> {
                    togglingDeviceIds = emptySet()
                    devicesState = parseDevices(body)
                }
                DevicesErrorPath -> {
                    togglingDeviceIds = emptySet()
                    devicesState = DevicesScreenState.Error(errorMessage(body, "Devices unavailable"))
                }
            }
        }
    }

    private fun requestWeather() {
        weatherState = WeatherScreenState.Connecting
        sendToPhone(
            path = WeatherRequestPath,
            payload = "{}",
            onSending = { weatherState = WeatherScreenState.Loading },
            onNoPhone = { weatherState = WeatherScreenState.Disconnected },
            onFailure = { weatherState = WeatherScreenState.Error(it) },
        )
    }

    private fun requestDevices() {
        devicesState = DevicesScreenState.Loading
        sendToPhone(
            path = DevicesRequestPath,
            payload = "{}",
            onSending = { devicesState = DevicesScreenState.Loading },
            onNoPhone = { devicesState = DevicesScreenState.Disconnected },
            onFailure = { devicesState = DevicesScreenState.Error(it) },
        )
    }

    private fun toggleDevice(device: WatchDevice) {
        if (!device.canToggle || togglingDeviceIds.contains(device.id)) return

        togglingDeviceIds = togglingDeviceIds + device.id
        val targetState = if (device.state.lowercase() == "on") "off" else "on"
        val payload = JSONObject()
            .put("houseId", device.houseId)
            .put("deviceId", device.id)
            .put("targetState", targetState)
            .toString()

        sendToPhone(
            path = DeviceTogglePath,
            payload = payload,
            onSending = {},
            onNoPhone = {
                togglingDeviceIds = togglingDeviceIds - device.id
                devicesState = DevicesScreenState.Disconnected
            },
            onFailure = {
                togglingDeviceIds = togglingDeviceIds - device.id
                devicesState = DevicesScreenState.Error(it)
            },
        )
    }

    private fun sendToPhone(
        path: String,
        payload: String,
        onSending: () -> Unit,
        onNoPhone: () -> Unit,
        onFailure: (String) -> Unit,
    ) {
        Wearable.getCapabilityClient(this)
            .getCapability(PhoneCapability, CapabilityClient.FILTER_REACHABLE)
            .addOnSuccessListener { capabilityInfo ->
                val node = capabilityInfo.nodes.firstOrNull()
                if (node == null) {
                    onNoPhone()
                    return@addOnSuccessListener
                }

                onSending()
                Wearable.getMessageClient(this)
                    .sendMessage(node.id, path, payload.toByteArray(StandardCharsets.UTF_8))
                    .addOnFailureListener { error ->
                        onFailure(error.message ?: "Could not send command")
                    }
            }
            .addOnFailureListener { error ->
                onFailure(error.message ?: "Could not find phone")
            }
    }

    private fun parseWeather(body: String): WeatherScreenState {
        return runCatching {
            val weather = JSONObject(body).getJSONObject("weather")
            WeatherScreenState.Weather(
                temperatureC = weather.getInt("temperatureC"),
                temperatureLabel = weather.optString("temperatureLabel"),
                condition = weather.getString("condition"),
                updatedAt = weather.optString("updatedAt"),
            )
        }.getOrElse { error ->
            WeatherScreenState.Error(error.message ?: "Invalid weather response")
        }
    }

    private fun parseDevices(body: String): DevicesScreenState {
        return runCatching {
            val devicesJson = JSONObject(body).getJSONArray("devices")
            val devices = buildList {
                for (index in 0 until devicesJson.length()) {
                    val device = devicesJson.getJSONObject(index)
                    add(
                        WatchDevice(
                            id = device.optString("id"),
                            houseId = device.optString("houseId"),
                            name = device.optString("name"),
                            state = device.optString("state", "unknown"),
                            availability = device.optString("availability", "unknown"),
                            brightness = if (device.has("brightness") && !device.isNull("brightness")) {
                                device.optInt("brightness").coerceIn(0, 100)
                            } else {
                                null
                            },
                        ),
                    )
                }
            }.filter { it.id.isNotBlank() && it.houseId.isNotBlank() }
            DevicesScreenState.Devices(devices)
        }.getOrElse { error ->
            DevicesScreenState.Error(error.message ?: "Invalid devices response")
        }
    }

    private fun errorMessage(body: String, fallback: String): String {
        return runCatching { JSONObject(body).optString("message") }
            .getOrNull()
            ?.takeIf { it.isNotBlank() }
            ?: fallback
    }
}

private sealed interface WeatherScreenState {
    data object Connecting : WeatherScreenState
    data object Loading : WeatherScreenState
    data object Disconnected : WeatherScreenState
    data class Weather(
        val temperatureC: Int,
        val temperatureLabel: String,
        val condition: String,
        val updatedAt: String,
    ) : WeatherScreenState

    data class Error(val message: String) : WeatherScreenState
}

private sealed interface DevicesScreenState {
    data object Loading : DevicesScreenState
    data object Disconnected : DevicesScreenState
    data class Devices(val devices: List<WatchDevice>) : DevicesScreenState
    data class Error(val message: String) : DevicesScreenState
}

private data class WatchDevice(
    val id: String,
    val houseId: String,
    val name: String,
    val state: String,
    val availability: String,
    val brightness: Int?,
) {
    val isOffline: Boolean get() = availability.lowercase() == "offline"
    val canToggle: Boolean get() = !isOffline && state.lowercase() in setOf("on", "off")
    val displayStatus: String
        get() = if (isOffline) {
            "offline"
        } else {
            when (state.lowercase()) {
                "on" -> "on"
                "off" -> "off"
                else -> "unknown"
            }
        }
}

@Composable
private fun HataWatchApp(
    weatherState: WeatherScreenState,
    devicesState: DevicesScreenState,
    togglingDeviceIds: Set<String>,
    onRetryWeather: () -> Unit,
    onRetryDevices: () -> Unit,
    onDeviceTap: (WatchDevice) -> Unit,
) {
    MaterialTheme {
        val pagerState = rememberPagerState(pageCount = { 2 })
        LaunchedEffect(pagerState.currentPage) {
            if (pagerState.currentPage == 1 && devicesState is DevicesScreenState.Loading) {
                onRetryDevices()
            }
        }

        Box(modifier = Modifier.fillMaxSize()) {
            HorizontalPager(
                state = pagerState,
                modifier = Modifier.fillMaxSize(),
            ) { page ->
                when (page) {
                    0 -> WeatherPage(
                        state = weatherState,
                        onRetry = onRetryWeather,
                    )
                    1 -> DevicesPage(
                        state = devicesState,
                        togglingDeviceIds = togglingDeviceIds,
                        onRetry = onRetryDevices,
                        onDeviceTap = onDeviceTap,
                    )
                }
            }
        }
    }
}

@Composable
private fun WeatherPage(
    state: WeatherScreenState,
    onRetry: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Box(
        modifier = modifier.fillMaxSize(),
        contentAlignment = Alignment.Center,
    ) {
        Column(
            modifier = Modifier.padding(horizontal = 16.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center,
        ) {
            when (state) {
                WeatherScreenState.Connecting -> StatusText("Connecting...")
                WeatherScreenState.Loading -> StatusText("Loading weather...")
                WeatherScreenState.Disconnected -> RetryState(
                    title = "Open Hata on your phone",
                    onRetry = onRetry,
                )
                is WeatherScreenState.Error -> RetryState(
                    title = state.message,
                    onRetry = onRetry,
                )
                is WeatherScreenState.Weather -> WeatherState(state)
            }
        }
    }
}

@OptIn(ExperimentalMaterialApi::class)
@Composable
private fun DevicesPage(
    state: DevicesScreenState,
    togglingDeviceIds: Set<String>,
    onRetry: () -> Unit,
    onDeviceTap: (WatchDevice) -> Unit,
    modifier: Modifier = Modifier,
) {
    val isRefreshing = state is DevicesScreenState.Loading
    val refreshState = rememberPullRefreshState(
        refreshing = isRefreshing,
        onRefresh = onRetry,
    )

    Box(
        modifier = modifier
            .fillMaxSize()
            .pullRefresh(refreshState),
        contentAlignment = Alignment.Center,
    ) {
        when (state) {
            DevicesScreenState.Loading -> LoadingState("Loading devices...")
            DevicesScreenState.Disconnected -> RetryState(
                title = "Open Hata on your phone",
                onRetry = onRetry,
            )
            is DevicesScreenState.Error -> RetryState(
                title = state.message,
                onRetry = onRetry,
            )
            is DevicesScreenState.Devices -> DeviceList(
                devices = state.devices,
                togglingDeviceIds = togglingDeviceIds,
                onRetry = onRetry,
                onDeviceTap = onDeviceTap,
            )
        }

        PullRefreshIndicator(
            refreshing = isRefreshing,
            state = refreshState,
            modifier = Modifier.align(Alignment.TopCenter),
        )
    }
}

@Composable
private fun DeviceList(
    devices: List<WatchDevice>,
    togglingDeviceIds: Set<String>,
    onRetry: () -> Unit,
    onDeviceTap: (WatchDevice) -> Unit,
) {
    if (devices.isEmpty()) {
        RetryState(title = "No devices", onRetry = onRetry)
        return
    }

    val listState = rememberScalingLazyListState()
    Box(modifier = Modifier.fillMaxSize()) {
        ScalingLazyColumn(
            modifier = Modifier.fillMaxSize(),
            state = listState,
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            item {
                Text(
                    text = "Devices",
                    style = MaterialTheme.typography.title3,
                    textAlign = TextAlign.Center,
                    modifier = Modifier.padding(bottom = 4.dp),
                )
            }
            items(devices.size) { index ->
                DeviceRow(
                    device = devices[index],
                    isToggling = togglingDeviceIds.contains(devices[index].id),
                    onClick = { onDeviceTap(devices[index]) },
                )
            }
        }
        PositionIndicator(
            scalingLazyListState = listState,
            modifier = Modifier.align(Alignment.CenterEnd),
        )
    }
}

@Composable
private fun DeviceRow(
    device: WatchDevice,
    isToggling: Boolean,
    onClick: () -> Unit,
) {
    val brightness = device.brightness ?: 0
    val fillFraction = (brightness.coerceIn(0, 100) / 100f).coerceIn(0f, 1f)
    val isOn = device.displayStatus == "on"
    val fillColor = if (isOn) Color(0xFF3D7DFF) else Color(0xFF303030)
    val backgroundColor = if (device.isOffline) Color(0xFF191919) else Color(0xFF242424)
    val borderColor = when {
        device.isOffline -> Color(0xCCFB4934)
        isOn -> Color(0xFFFABD2F)
        else -> Color(0xFF504945)
    }
    val shape = RoundedCornerShape(18.dp)

    Box(
        modifier = Modifier
            .fillMaxWidth(0.86f)
            .height(54.dp)
            .clip(shape)
            .background(backgroundColor)
            .border(width = 1.dp, color = borderColor, shape = shape)
            .clickable(enabled = device.canToggle && !isToggling, onClick = onClick),
    ) {
        Box(
            modifier = Modifier
                .fillMaxHeight()
                .fillMaxWidth(fillFraction)
                .background(fillColor),
        )
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(horizontal = 14.dp, vertical = 7.dp),
            verticalArrangement = Arrangement.Center,
        ) {
            Text(
                text = device.name.ifBlank { device.id },
                maxLines = 1,
                overflow = TextOverflow.Ellipsis,
                style = MaterialTheme.typography.title3,
            )
            Text(
                text = deviceStatusLabel(device),
                maxLines = 1,
                overflow = TextOverflow.Ellipsis,
                style = MaterialTheme.typography.caption2,
            )
        }
        if (isToggling) {
            CircularProgressIndicator(
                modifier = Modifier
                    .align(Alignment.CenterEnd)
                    .padding(end = 12.dp)
                    .width(18.dp)
                    .height(18.dp),
                strokeWidth = 2.dp,
            )
        }
    }
}

private fun deviceStatusLabel(device: WatchDevice): String {
    val brightness = device.brightness
    val brightnessLabel = if (brightness == null) "" else " • $brightness%"
    return device.displayStatus + brightnessLabel
}

@Composable
private fun LoadingState(text: String) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        CircularProgressIndicator()
        Spacer(modifier = Modifier.height(8.dp))
        StatusText(text)
    }
}

@Composable
private fun StatusText(text: String) {
    Text(
        text = text,
        textAlign = TextAlign.Center,
        style = MaterialTheme.typography.title3,
    )
}

@Composable
private fun RetryState(
    title: String,
    onRetry: () -> Unit,
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
        modifier = Modifier.padding(horizontal = 16.dp),
    ) {
        StatusText(title)
        Spacer(modifier = Modifier.height(8.dp))
        Button(onClick = onRetry) {
            Text("Retry")
        }
    }
}

@Composable
private fun WeatherState(state: WeatherScreenState.Weather) {
    Text(
        text = state.temperatureLabel.ifBlank { "${state.temperatureC}°C" },
        textAlign = TextAlign.Center,
        style = MaterialTheme.typography.display2,
    )
    Text(
        text = state.condition,
        textAlign = TextAlign.Center,
        style = MaterialTheme.typography.title3,
    )
    val updatedLabel = updatedLabel(state.updatedAt)
    if (updatedLabel != null) {
        Text(
            text = updatedLabel,
            textAlign = TextAlign.Center,
            style = MaterialTheme.typography.caption2,
        )
    }
}

private fun updatedLabel(updatedAt: String): String? {
    if (updatedAt.isBlank()) return null

    val instant = runCatching { Instant.parse(updatedAt) }.getOrNull() ?: return null
    var age = Duration.between(instant, Instant.now())
    if (age.isNegative) age = Duration.ZERO

    val minutes = age.toMinutes()
    if (minutes < 1) return "updated just now"
    if (minutes < 60) {
        return if (minutes == 1L) {
            "updated 1 minute ago"
        } else {
            "updated $minutes minutes ago"
        }
    }

    val hours = age.toHours()
    if (hours < 24) {
        return if (hours == 1L) {
            "updated 1 hour ago"
        } else {
            "updated $hours hours ago"
        }
    }

    val days = hours / 24
    return if (days == 1L) {
        "updated 1 day ago"
    } else {
        "updated $days days ago"
    }
}

@Preview(
    showBackground = true,
    backgroundColor = 0xFF000000,
)
@Composable
private fun HataWatchAppPreview() {
    HataWatchApp(
        weatherState = WeatherScreenState.Weather(
            temperatureC = 22,
            temperatureLabel = "22°C",
            condition = "Cloudy",
            updatedAt = "2026-08-19T20:00:00Z",
        ),
        devicesState = DevicesScreenState.Devices(
            listOf(
                WatchDevice("1", "h1", "Kitchen", "on", "online", 80),
                WatchDevice("2", "h1", "Hall", "off", "online", 0),
                WatchDevice("3", "h1", "Porch", "off", "offline", null),
            ),
        ),
        togglingDeviceIds = setOf("2"),
        onRetryWeather = {},
        onRetryDevices = {},
        onDeviceTap = {},
    )
}
