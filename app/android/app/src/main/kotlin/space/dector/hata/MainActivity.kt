package space.dector.hata

import android.os.Bundle
import com.google.android.gms.wearable.MessageClient
import com.google.android.gms.wearable.MessageEvent
import com.google.android.gms.wearable.Wearable
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import org.json.JSONArray
import org.json.JSONObject
import java.nio.charset.StandardCharsets

private const val WearChannelName = "hata/wear"
private const val WeatherRequestPath = "/watch/weather/request"
private const val WeatherUpdatePath = "/phone/weather/update"
private const val WeatherErrorPath = "/phone/weather/error"
private const val DevicesRequestPath = "/watch/devices/request"
private const val DeviceTogglePath = "/watch/device/toggle"
private const val DevicesUpdatePath = "/phone/devices/update"
private const val DevicesErrorPath = "/phone/devices/error"

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
        when (event.path) {
            WeatherRequestPath -> handleWeatherRequest(event)
            DevicesRequestPath -> handleDevicesRequest(event)
            DeviceTogglePath -> handleDeviceToggle(event)
        }
    }

    private fun handleWeatherRequest(event: MessageEvent) {
        invokeFlutter(
            nodeId = event.sourceNodeId,
            method = "getWeather",
            arguments = null,
            updatePath = WeatherUpdatePath,
            errorPath = WeatherErrorPath,
            notImplementedMessage = "Weather bridge is not implemented",
            payload = ::weatherPayload,
        )
    }

    private fun handleDevicesRequest(event: MessageEvent) {
        invokeFlutter(
            nodeId = event.sourceNodeId,
            method = "getDevices",
            arguments = null,
            updatePath = DevicesUpdatePath,
            errorPath = DevicesErrorPath,
            notImplementedMessage = "Devices bridge is not implemented",
            payload = ::devicesPayload,
        )
    }

    private fun handleDeviceToggle(event: MessageEvent) {
        val body = event.data.toString(StandardCharsets.UTF_8)
        val json = runCatching { JSONObject(body) }.getOrNull()
        if (json == null) {
            sendError(event.sourceNodeId, DevicesErrorPath, "Invalid toggle payload")
            return
        }

        val arguments = hashMapOf<String, Any?>(
            "houseId" to json.optString("houseId"),
            "deviceId" to json.optString("deviceId"),
            "targetState" to json.optString("targetState"),
        )
        invokeFlutter(
            nodeId = event.sourceNodeId,
            method = "toggleDevice",
            arguments = arguments,
            updatePath = DevicesUpdatePath,
            errorPath = DevicesErrorPath,
            notImplementedMessage = "Device controls bridge is not implemented",
            payload = ::devicesPayload,
        )
    }

    private fun invokeFlutter(
        nodeId: String,
        method: String,
        arguments: Any?,
        updatePath: String,
        errorPath: String,
        notImplementedMessage: String,
        payload: (Any?) -> String,
    ) {
        runOnUiThread {
            val channel = wearChannel
            if (channel == null) {
                sendError(nodeId, errorPath, "Phone app is not ready")
                return@runOnUiThread
            }

            channel.invokeMethod(
                method,
                arguments,
                object : MethodChannel.Result {
                    override fun success(result: Any?) {
                        sendMessage(nodeId, updatePath, payload(result))
                    }

                    override fun error(
                        errorCode: String,
                        errorMessage: String?,
                        errorDetails: Any?,
                    ) {
                        sendError(nodeId, errorPath, errorMessage ?: errorCode)
                    }

                    override fun notImplemented() {
                        sendError(nodeId, errorPath, notImplementedMessage)
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

    private fun devicesPayload(result: Any?): String {
        val devices = JSONArray()
        if (result is List<*>) {
            result.forEach { item ->
                if (item is Map<*, *>) {
                    devices.put(
                        JSONObject()
                            .put("id", item["id"])
                            .put("houseId", item["houseId"])
                            .put("name", item["name"])
                            .put("state", item["state"])
                            .put("availability", item["availability"])
                            .put("brightness", item["brightness"]),
                    )
                }
            }
        }
        return JSONObject()
            .put("devices", devices)
            .toString()
    }

    private fun sendError(
        nodeId: String,
        path: String,
        message: String,
    ) {
        val payload = JSONObject()
            .put("message", message)
            .toString()
        sendMessage(nodeId, path, payload)
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
