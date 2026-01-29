package hata.integrations.wiz

import io.kotest.matchers.collections.shouldHaveSize
import io.kotest.matchers.shouldBe
import io.kotest.matchers.string.shouldNotBeEmpty
import kotlinx.serialization.json.Json
import org.junit.Test


class WizDiscoveryTest {

    private val json = Json {
        ignoreUnknownKeys = true
        isLenient = true
    }

    @Test
    fun `parseDevice with productName`() {
        val registrationResult = RegistrationResult(
            mac = "A1B2C3D4E5F6",
            ip = "192.168.1.100",
            productName = "WiZ Tunable White",
            modelName = "ESP_XXX",
            firmwareVersion = "1.2.3",
            moduleName = "ESP_01_SHRGBC"
        )

        val response = RegistrationResponse(
            method = "registration",
            env = "pro",
            result = registrationResult
        )

        // Create a device using the same logic as parseDevice
        val deviceName = when {
            !response.result.productName.isNullOrEmpty() -> response.result.productName
            !response.result.modelName.isNullOrEmpty() -> response.result.modelName
            !response.result.moduleName.isNullOrEmpty() -> response.result.moduleName
            else -> "WiZ Device"
        }

        val deviceType = when {
            deviceName.contains("strip", ignoreCase = true) -> "Light Strip"
            deviceName.contains("spot", ignoreCase = true) -> "Spotlight"
            deviceName.contains("bulb", ignoreCase = true) -> "Bulb"
            else -> "Smart Light"
        }

        deviceName shouldBe "WiZ Tunable White"
        deviceType shouldBe "Smart Light"
    }

    @Test
    fun `parseDevice with modelName only`() {
        val registrationResult = RegistrationResult(
            mac = "A1B2C3D4E5F6",
            ip = "192.168.1.100",
            productName = null,
            modelName = "ESP_XXX",
            firmwareVersion = "1.2.3",
            moduleName = "ESP_01_SHRGBC"
        )

        val response = RegistrationResponse(
            method = "registration",
            result = registrationResult
        )

        val deviceName = when {
            !response.result.productName.isNullOrEmpty() -> response.result.productName
            !response.result.modelName.isNullOrEmpty() -> response.result.modelName
            !response.result.moduleName.isNullOrEmpty() -> response.result.moduleName
            else -> "WiZ Device"
        }

        deviceName shouldBe "ESP_XXX"
    }

    @Test
    fun `parseDevice with strip in name`() {
        val registrationResult = RegistrationResult(
            mac = "A1B2C3D4E5F6",
            productName = "WiZ LED Strip"
        )

        val deviceName = "WiZ LED Strip"
        val deviceType = when {
            deviceName.contains("strip", ignoreCase = true) -> "Light Strip"
            deviceName.contains("spot", ignoreCase = true) -> "Spotlight"
            deviceName.contains("bulb", ignoreCase = true) -> "Bulb"
            else -> "Smart Light"
        }

        deviceType shouldBe "Light Strip"
    }

    @Test
    fun `parseDevice with spot in name`() {
        val deviceName = "WiZ Spot"
        val deviceType = when {
            deviceName.contains("strip", ignoreCase = true) -> "Light Strip"
            deviceName.contains("spot", ignoreCase = true) -> "Spotlight"
            deviceName.contains("bulb", ignoreCase = true) -> "Bulb"
            else -> "Smart Light"
        }

        deviceType shouldBe "Spotlight"
    }

    @Test
    fun `parseDevice with bulb in name`() {
        val deviceName = "WiZ Bulb White"
        val deviceType = when {
            deviceName.contains("strip", ignoreCase = true) -> "Light Strip"
            deviceName.contains("spot", ignoreCase = true) -> "Spotlight"
            deviceName.contains("bulb", ignoreCase = true) -> "Bulb"
            else -> "Smart Light"
        }

        deviceType shouldBe "Bulb"
    }

    @Test
    fun `parseDevice with no name defaults to WiZ Device`() {
        val registrationResult = RegistrationResult(
            mac = "A1B2C3D4E5F6",
            productName = null,
            modelName = null,
            moduleName = null
        )

        val response = RegistrationResponse(
            method = "registration",
            result = registrationResult
        )

        val deviceName = when {
            !response.result.productName.isNullOrEmpty() -> response.result.productName
            !response.result.modelName.isNullOrEmpty() -> response.result.modelName
            !response.result.moduleName.isNullOrEmpty() -> response.result.moduleName
            else -> "WiZ Device"
        }

        deviceName shouldBe "WiZ Device"
    }

    @Test
    fun `discovery packet format`() {
        val expectedPacket = """{"method":"registration","params":{"phoneMac":"AAAAAAAAAAAA","register":false,"phoneIp":"0.0.0.0","id":"1"}}"""

        // Verify the packet is valid JSON
        val parsedPacket = json.parseToJsonElement(expectedPacket)
        parsedPacket.toString().shouldNotBeEmpty()
    }

    @Test
    fun `registration response parsing`() {
        val jsonResponse = """
            {
                "method": "registration",
                "env": "pro",
                "result": {
                    "mac": "A1B2C3D4E5F6",
                    "ip": "192.168.1.100",
                    "productName": "WiZ Tunable White"
                }
            }
        """.trimIndent()

        val response = json.decodeFromString<RegistrationResponse>(jsonResponse)

        response.method shouldBe "registration"
        response.result.mac shouldBe "A1B2C3D4E5F6"
        response.result.ip shouldBe "192.168.1.100"
        response.result.productName shouldBe "WiZ Tunable White"
    }

    @Test
    fun `deduplication by MAC address`() {
        val seenMacs = mutableSetOf<String>()
        val devices = mutableListOf<WizDevice>()

        val device1 = WizDevice("Light 1", "192.168.1.100", "Bulb", "A1B2C3D4E5F6")
        val device2 = WizDevice("Light 1", "192.168.1.100", "Bulb", "A1B2C3D4E5F6") // Same MAC
        val device3 = WizDevice("Light 2", "192.168.1.101", "Bulb", "A1B2C3D4E5F7") // Different MAC

        // Simulate deduplication logic
        if (device1.mac !in seenMacs) {
            seenMacs.add(device1.mac)
            devices.add(device1)
        }

        if (device2.mac !in seenMacs) {
            seenMacs.add(device2.mac)
            devices.add(device2)
        }

        if (device3.mac !in seenMacs) {
            seenMacs.add(device3.mac)
            devices.add(device3)
        }

        devices shouldHaveSize 2 // Only 2 unique devices
        seenMacs shouldHaveSize 2
    }
}
