package space.dector.hata

import android.os.Bundle
import com.google.android.gms.wearable.MessageClient
import com.google.android.gms.wearable.MessageEvent
import com.google.android.gms.wearable.Wearable
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import org.json.JSONObject
import java.nio.charset.StandardCharsets

private const val WearChannelName = "hata/wear"
private const val WeatherRequestPath = "/watch/weather/request"
private const val WeatherUpdatePath = "/phone/weather/update"
private const val WeatherErrorPath = "/phone/weather/error"

class MainActivity : FlutterActivity(), MessageClient.OnMessageReceivedListener {
    private var wearChannel: MethodChannel? = null

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)
        wearChannel = MethodChannel(flutterEngine.dartExecutor.binaryMessenger, WearChannelName)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Wearable.getMessageClient(this).addListener(this)
    }

    override fun onDestroy() {
        Wearable.getMessageClient(this).removeListener(this)
        wearChannel = null
        super.onDestroy()
    }

    override fun onMessageReceived(event: MessageEvent) {
        if (event.path != WeatherRequestPath) {
            return
        }

        runOnUiThread {
            val channel = wearChannel
            if (channel == null) {
                sendError(event.sourceNodeId, "Phone app is not ready")
                return@runOnUiThread
            }

            channel.invokeMethod(
                "getWeather",
                null,
                object : MethodChannel.Result {
                    override fun success(result: Any?) {
                        val payload = weatherPayload(result)
                        sendMessage(event.sourceNodeId, WeatherUpdatePath, payload)
                    }

                    override fun error(
                        errorCode: String,
                        errorMessage: String?,
                        errorDetails: Any?,
                    ) {
                        sendError(event.sourceNodeId, errorMessage ?: errorCode)
                    }

                    override fun notImplemented() {
                        sendError(event.sourceNodeId, "Weather bridge is not implemented")
                    }
                },
            )
        }
    }

    private fun weatherPayload(result: Any?): String {
        val weather = JSONObject()
        if (result is Map<*, *>) {
            weather.put("temperatureC", result["temperatureC"])
            weather.put("temperatureLabel", result["temperatureLabel"])
            weather.put("condition", result["condition"])
            weather.put("updatedAt", result["updatedAt"])
        }
        return JSONObject()
            .put("weather", weather)
            .toString()
    }

    private fun sendError(
        nodeId: String,
        message: String,
    ) {
        val payload = JSONObject()
            .put("message", message)
            .toString()
        sendMessage(nodeId, WeatherErrorPath, payload)
    }

    private fun sendMessage(
        nodeId: String,
        path: String,
        payload: String,
    ) {
        Wearable.getMessageClient(this)
            .sendMessage(nodeId, path, payload.toByteArray(StandardCharsets.UTF_8))
    }
}
