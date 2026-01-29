package hata.integrations.wiz

import io.kotest.matchers.shouldBe
import io.kotest.matchers.nulls.shouldBeNull
import kotlinx.serialization.json.Json
import org.junit.Test


class WizModelsTest {

    private val json = Json {
        ignoreUnknownKeys = true
        isLenient = true
    }

    @Test
    fun `WizDevice serialization and deserialization`() {
        val device = WizDevice(
            name = "Test Light",
            ip = "192.168.1.100",
            type = "Bulb",
            mac = "A1:B2:C3:D4:E5:F6"
        )

        val jsonString = json.encodeToString(device)
        val deserializedDevice = json.decodeFromString<WizDevice>(jsonString)

        deserializedDevice.name shouldBe device.name
        deserializedDevice.ip shouldBe device.ip
        deserializedDevice.type shouldBe device.type
        deserializedDevice.mac shouldBe device.mac
    }

    @Test
    fun `PilotState serialization and deserialization`() {
        val pilotState = PilotState(
            state = true,
            dimming = 75,
            r = 255,
            g = 128,
            b = 64,
            temp = 3000,
            sceneId = 1
        )

        val jsonString = json.encodeToString(pilotState)
        val deserializedState = json.decodeFromString<PilotState>(jsonString)

        deserializedState.state shouldBe pilotState.state
        deserializedState.dimming shouldBe pilotState.dimming
        deserializedState.r shouldBe pilotState.r
        deserializedState.g shouldBe pilotState.g
        deserializedState.b shouldBe pilotState.b
        deserializedState.temp shouldBe pilotState.temp
        deserializedState.sceneId shouldBe pilotState.sceneId
    }

    @Test
    fun `RegistrationResponse deserialization from JSON`() {
        val jsonResponse = """
            {
                "method": "registration",
                "env": "pro",
                "result": {
                    "mac": "A1B2C3D4E5F6",
                    "ip": "192.168.1.100",
                    "productName": "WiZ Tunable White",
                    "modelName": "ESP_XXX",
                    "fwVersion": "1.2.3",
                    "moduleName": "ESP_01_SHRGBC"
                }
            }
        """.trimIndent()

        val response = json.decodeFromString<RegistrationResponse>(jsonResponse)

        response.method shouldBe "registration"
        response.env shouldBe "pro"
        response.result.mac shouldBe "A1B2C3D4E5F6"
        response.result.ip shouldBe "192.168.1.100"
        response.result.productName shouldBe "WiZ Tunable White"
        response.result.modelName shouldBe "ESP_XXX"
        response.result.firmwareVersion shouldBe "1.2.3"
        response.result.moduleName shouldBe "ESP_01_SHRGBC"
    }

    @Test
    fun `RegistrationResponse with minimal fields`() {
        val jsonResponse = """
            {
                "method": "registration",
                "result": {
                    "mac": "A1B2C3D4E5F6"
                }
            }
        """.trimIndent()

        val response = json.decodeFromString<RegistrationResponse>(jsonResponse)

        response.method shouldBe "registration"
        response.result.mac shouldBe "A1B2C3D4E5F6"
        response.result.productName.shouldBeNull()
        response.result.modelName.shouldBeNull()
    }

    @Test
    fun `GetPilotResponse deserialization from JSON`() {
        val jsonResponse = """
            {
                "method": "getPilot",
                "env": "pro",
                "result": {
                    "state": true,
                    "dimming": 50,
                    "r": 255,
                    "g": 0,
                    "b": 128,
                    "temp": 3200,
                    "sceneId": 0
                }
            }
        """.trimIndent()

        val response = json.decodeFromString<GetPilotResponse>(jsonResponse)

        response.method shouldBe "getPilot"
        response.env shouldBe "pro"
        response.result.state shouldBe true
        response.result.dimming shouldBe 50
        response.result.r shouldBe 255
        response.result.g shouldBe 0
        response.result.b shouldBe 128
        response.result.temp shouldBe 3200
        response.result.sceneId shouldBe 0
    }

    @Test
    fun `PilotState default values`() {
        val pilotState = PilotState(state = true)

        pilotState.state shouldBe true
        pilotState.dimming shouldBe 0
        pilotState.r shouldBe 0
        pilotState.g shouldBe 0
        pilotState.b shouldBe 0
        pilotState.temp shouldBe 0
        pilotState.sceneId shouldBe 0
    }

    @Test
    fun `Result Success type`() {
        val result = Result.Success(42)

        result.data shouldBe 42
    }

    @Test
    fun `Result Error type`() {
        val exception = IllegalArgumentException("Test error")
        val result = Result.Error(exception)

        result.exception shouldBe exception
        result.exception.message shouldBe "Test error"
    }
}
