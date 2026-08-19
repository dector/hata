package space.dector.hata.watch

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.wear.compose.material.Button
import androidx.wear.compose.material.MaterialTheme
import androidx.wear.compose.material.Text
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

class MainActivity : ComponentActivity(), MessageClient.OnMessageReceivedListener {
    private var screenState by mutableStateOf<WatchScreenState>(WatchScreenState.Connecting)

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        setContent {
            HataWatchApp(
                state = screenState,
                onRetry = ::requestWeather,
            )
        }
    }

    override fun onResume() {
        super.onResume()
        Wearable.getMessageClient(this).addListener(this)
        requestWeather()
    }

    override fun onPause() {
        Wearable.getMessageClient(this).removeListener(this)
        super.onPause()
    }

    override fun onMessageReceived(event: MessageEvent) {
        val body = event.data.toString(StandardCharsets.UTF_8)
        runOnUiThread {
            when (event.path) {
                WeatherUpdatePath -> screenState = parseWeather(body)
                WeatherErrorPath -> screenState = WatchScreenState.Error(
                    message = runCatching { JSONObject(body).optString("message") }
                        .getOrNull()
                        ?.takeIf { it.isNotBlank() }
                        ?: "Weather unavailable",
                )
            }
        }
    }

    private fun requestWeather() {
        screenState = WatchScreenState.Connecting

        Wearable.getCapabilityClient(this)
            .getCapability(PhoneCapability, CapabilityClient.FILTER_REACHABLE)
            .addOnSuccessListener { capabilityInfo ->
                val node = capabilityInfo.nodes.firstOrNull()
                if (node == null) {
                    screenState = WatchScreenState.Disconnected
                    return@addOnSuccessListener
                }

                screenState = WatchScreenState.Loading
                Wearable.getMessageClient(this)
                    .sendMessage(node.id, WeatherRequestPath, "{}".toByteArray(StandardCharsets.UTF_8))
                    .addOnFailureListener { error ->
                        screenState = WatchScreenState.Error(error.message ?: "Could not request weather")
                    }
            }
            .addOnFailureListener { error ->
                screenState = WatchScreenState.Error(error.message ?: "Could not find phone")
            }
    }

    private fun parseWeather(body: String): WatchScreenState {
        return runCatching {
            val weather = JSONObject(body).getJSONObject("weather")
            WatchScreenState.Weather(
                temperatureC = weather.getInt("temperatureC"),
                temperatureLabel = weather.optString("temperatureLabel"),
                condition = weather.getString("condition"),
                updatedAt = weather.optString("updatedAt"),
            )
        }.getOrElse { error ->
            WatchScreenState.Error(error.message ?: "Invalid weather response")
        }
    }
}

private sealed interface WatchScreenState {
    data object Connecting : WatchScreenState
    data object Loading : WatchScreenState
    data object Disconnected : WatchScreenState
    data class Weather(
        val temperatureC: Int,
        val temperatureLabel: String,
        val condition: String,
        val updatedAt: String,
    ) : WatchScreenState

    data class Error(val message: String) : WatchScreenState
}

@Composable
private fun HataWatchApp(
    state: WatchScreenState,
    onRetry: () -> Unit,
) {
    MaterialTheme {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center,
        ) {
            WatchContent(
                state = state,
                onRetry = onRetry,
            )
        }
    }
}

@Composable
private fun WatchContent(
    state: WatchScreenState,
    onRetry: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Column(
        modifier = modifier.padding(horizontal = 16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        when (state) {
            WatchScreenState.Connecting -> StatusText("Connecting...")
            WatchScreenState.Loading -> StatusText("Loading weather...")
            WatchScreenState.Disconnected -> RetryState(
                title = "Open Hata on your phone",
                onRetry = onRetry,
            )
            is WatchScreenState.Error -> RetryState(
                title = state.message,
                onRetry = onRetry,
            )
            is WatchScreenState.Weather -> WeatherState(state)
        }
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
    StatusText(title)
    Spacer(modifier = Modifier.height(8.dp))
    Button(onClick = onRetry) {
        Text("Retry")
    }
}

@Composable
private fun WeatherState(state: WatchScreenState.Weather) {
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
        state = WatchScreenState.Weather(
            temperatureC = 22,
            temperatureLabel = "22°C",
            condition = "Cloudy",
            updatedAt = "2026-08-19T20:00:00Z",
        ),
        onRetry = {},
    )
}
