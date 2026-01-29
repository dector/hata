package hata.integrations.wiz

import android.util.Log
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.async
import kotlinx.coroutines.delay
import kotlinx.coroutines.isActive
import kotlinx.coroutines.withContext
import kotlinx.coroutines.withTimeoutOrNull
import kotlinx.serialization.json.Json
import java.net.DatagramPacket
import java.net.DatagramSocket
import java.net.InetAddress
import java.net.InetSocketAddress
import kotlin.time.Duration
import kotlin.time.Duration.Companion.seconds


private const val TAG = "WizDiscovery"
private const val WIZ_PORT = 38899
private const val BROADCAST_ADDR = "255.255.255.255"
private const val DISCOVERY_PACKET = """{"method":"registration","params":{"phoneMac":"AAAAAAAAAAAA","register":false,"phoneIp":"0.0.0.0","id":"1"}}"""

object WizDiscovery {

    private val json = Json {
        ignoreUnknownKeys = true
        isLenient = true
    }

    /**
     * Scans the network for WiZ devices by sending UDP broadcasts.
     *
     * @param duration How long to scan for devices
     * @return List of discovered WiZ devices
     */
    suspend fun scan(duration: Duration = 10.seconds): List<WizDevice> = withContext(Dispatchers.IO) {
        Log.d(TAG, "Starting WiZ device discovery for $duration...")

        val devices = mutableListOf<WizDevice>()
        val seenMacs = mutableSetOf<String>()

        val socket = DatagramSocket(null).apply {
            reuseAddress = true
            broadcast = true
            bind(InetSocketAddress(WIZ_PORT))
        }

        try {
            socket.use { sock ->
                // Start listening for responses
                val listenJob = async {
                    listenForResponses(sock, seenMacs, devices)
                }

                // Send broadcast packets periodically
                val broadcastJob = async {
                    broadcastDiscovery(sock)
                }

                // Wait for the specified duration
                withTimeoutOrNull(duration) {
                    listenJob.await()
                }

                // Cancel both jobs
                listenJob.cancel()
                broadcastJob.cancel()
            }
        } catch (e: Exception) {
            Log.e(TAG, "Error during WiZ discovery", e)
        }

        Log.d(TAG, "WiZ discovery completed, found ${devices.size} unique devices")
        devices.toList()
    }

    private suspend fun listenForResponses(
        socket: DatagramSocket,
        seenMacs: MutableSet<String>,
        devices: MutableList<WizDevice>
    ) = withContext(Dispatchers.IO) {
        val buffer = ByteArray(2048)
        val packet = DatagramPacket(buffer, buffer.size)

        socket.soTimeout = 1000 // 1 second timeout for reads

        while (isActive) {
            try {
                socket.receive(packet)

                val response = String(packet.data, 0, packet.length)
                val deviceIp = packet.address.hostAddress ?: continue

                try {
                    val wizResponse = json.decodeFromString<RegistrationResponse>(response)

                    if (wizResponse.result.mac.isNotEmpty() && wizResponse.result.mac !in seenMacs) {
                        seenMacs.add(wizResponse.result.mac)
                        val device = parseDevice(wizResponse, deviceIp)
                        devices.add(device)
                        Log.d(TAG, "Discovered device: ${device.name} at ${device.ip}")
                    }
                } catch (e: Exception) {
                    Log.w(TAG, "Failed to parse response from $deviceIp: ${e.message}")
                }
            } catch (e: java.net.SocketTimeoutException) {
                // Timeout is expected, continue listening
                continue
            } catch (e: Exception) {
                if (isActive) {
                    Log.w(TAG, "Error receiving packet: ${e.message}")
                }
            }
        }
    }

    private suspend fun broadcastDiscovery(socket: DatagramSocket) = withContext(Dispatchers.IO) {
        val broadcastAddr = InetAddress.getByName(BROADCAST_ADDR)
        val data = DISCOVERY_PACKET.toByteArray()
        val packet = DatagramPacket(data, data.size, broadcastAddr, WIZ_PORT)

        // Send first packet immediately
        try {
            socket.send(packet)
            Log.d(TAG, "Sent initial discovery packet")
        } catch (e: Exception) {
            Log.e(TAG, "Failed to send initial discovery packet", e)
        }

        // Send periodic packets
        while (isActive) {
            delay(1.seconds)
            try {
                socket.send(packet)
                Log.d(TAG, "Sent periodic discovery packet")
            } catch (e: Exception) {
                if (isActive) {
                    Log.e(TAG, "Failed to send periodic discovery packet", e)
                }
            }
        }
    }

    private fun parseDevice(response: RegistrationResponse, ip: String): WizDevice {
        val result = response.result

        val deviceName = when {
            !result.productName.isNullOrEmpty() -> result.productName
            !result.modelName.isNullOrEmpty() -> result.modelName
            !result.moduleName.isNullOrEmpty() -> result.moduleName
            else -> "WiZ Device"
        }

        val deviceType = when {
            deviceName.contains("strip", ignoreCase = true) -> "Light Strip"
            deviceName.contains("spot", ignoreCase = true) -> "Spotlight"
            deviceName.contains("bulb", ignoreCase = true) -> "Bulb"
            else -> "Smart Light"
        }

        return WizDevice(
            name = deviceName,
            ip = ip,
            type = deviceType,
            mac = result.mac
        )
    }
}
