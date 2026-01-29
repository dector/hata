package hata.integrations.wiz

import android.util.Log
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import java.net.DatagramPacket
import java.net.DatagramSocket
import java.net.InetAddress


private const val TAG = "WizControl"
private const val WIZ_PORT = 38899

class WizControl(private val device: WizDevice) {

    private val json = Json {
        ignoreUnknownKeys = true
        isLenient = true
    }

    /**
     * Turns the light on.
     */
    suspend fun turnOn(): Result<Unit> = sendPilotCommand(
        buildJsonObject { put("state", true) }
    )

    /**
     * Turns the light off.
     */
    suspend fun turnOff(): Result<Unit> = sendPilotCommand(
        buildJsonObject { put("state", false) }
    )

    /**
     * Toggles the light's power state.
     */
    suspend fun toggle(): Result<Unit> {
        return when (val stateResult = getState()) {
            is Result.Success -> {
                if (stateResult.data.state) {
                    turnOff()
                } else {
                    turnOn()
                }
            }
            is Result.Error -> Result.Error(stateResult.exception)
        }
    }

    /**
     * Retrieves the current state of the light.
     */
    suspend fun getState(): Result<PilotState> = withContext(Dispatchers.IO) {
        try {
            val command = buildJsonObject {
                put("method", "getPilot")
                put("params", buildJsonObject { })
            }

            val address = InetAddress.getByName(device.ip)
            val commandBytes = command.toString().toByteArray()

            DatagramSocket().use { socket ->
                socket.soTimeout = 2000 // 2 second timeout

                // Send command
                val sendPacket = DatagramPacket(
                    commandBytes,
                    commandBytes.size,
                    address,
                    WIZ_PORT
                )
                socket.send(sendPacket)

                // Receive response
                val buffer = ByteArray(2048)
                val receivePacket = DatagramPacket(buffer, buffer.size)
                socket.receive(receivePacket)

                val response = String(receivePacket.data, 0, receivePacket.length)

                // Parse response
                val pilotResponse = json.decodeFromString<GetPilotResponse>(response)

                if (pilotResponse.method != "getPilot") {
                    Result.Error(
                        IllegalStateException(
                            "Unexpected response method: ${pilotResponse.method}"
                        )
                    )
                } else {
                    Result.Success(pilotResponse.result)
                }
            }
        } catch (e: Exception) {
            Log.e(TAG, "Failed to get state for device '${device.name}' at ${device.ip}", e)
            Result.Error(e)
        }
    }

    /**
     * Sets the light's brightness.
     *
     * @param brightness Value between 10 and 100
     */
    suspend fun setBrightness(brightness: Int): Result<Unit> {
        if (brightness !in 10..100) {
            return Result.Error(
                IllegalArgumentException("Brightness must be between 10 and 100, got $brightness")
            )
        }

        return sendPilotCommand(
            buildJsonObject { put("dimming", brightness) }
        )
    }

    /**
     * Sets the light's RGB color.
     *
     * @param r Red value (0-255)
     * @param g Green value (0-255)
     * @param b Blue value (0-255)
     */
    suspend fun setRGBColor(r: Int, g: Int, b: Int): Result<Unit> {
        if (r !in 0..255 || g !in 0..255 || b !in 0..255) {
            return Result.Error(
                IllegalArgumentException("RGB values must be between 0 and 255")
            )
        }

        return sendPilotCommand(
            buildJsonObject {
                put("r", r)
                put("g", g)
                put("b", b)
            }
        )
    }

    /**
     * Sets the light's color temperature.
     *
     * @param temp Temperature in Kelvin (2200-6500)
     */
    suspend fun setColorTemperature(temp: Int): Result<Unit> {
        if (temp !in 2200..6500) {
            return Result.Error(
                IllegalArgumentException(
                    "Color temperature must be between 2200K and 6500K, got ${temp}K"
                )
            )
        }

        return sendPilotCommand(
            buildJsonObject { put("temp", temp) }
        )
    }

    /**
     * Sends a setPilot command to the light.
     */
    private suspend fun sendPilotCommand(params: JsonObject): Result<Unit> {
        val command = buildJsonObject {
            put("method", "setPilot")
            put("params", params)
        }

        return sendUDPCommand(command.toString())
    }

    /**
     * Sends a UDP command to the light without waiting for a response.
     */
    private suspend fun sendUDPCommand(commandJson: String): Result<Unit> = withContext(Dispatchers.IO) {
        try {
            val address = InetAddress.getByName(device.ip)
            val data = commandJson.toByteArray()
            val packet = DatagramPacket(data, data.size, address, WIZ_PORT)

            DatagramSocket().use { socket ->
                socket.send(packet)
                Log.d(TAG, "Sent command to ${device.name} at ${device.ip}: $commandJson")
            }

            Result.Success(Unit)
        } catch (e: Exception) {
            Log.e(TAG, "Failed to send command to device '${device.name}' at ${device.ip}", e)
            Result.Error(e)
        }
    }
}
