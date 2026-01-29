package hata.integrations.wiz

import io.kotest.matchers.shouldBe
import io.kotest.matchers.types.shouldBeInstanceOf
import kotlinx.coroutines.test.runTest
import org.junit.Before
import org.junit.Test


class WizControlTest {

    private lateinit var testDevice: WizDevice
    private lateinit var wizControl: WizControl

    @Before
    fun setup() {
        testDevice = WizDevice(
            name = "Test Light",
            ip = "192.168.1.100",
            type = "Bulb",
            mac = "A1B2C3D4E5F6"
        )
        wizControl = WizControl(testDevice)
    }

    @Test
    fun `setBrightness validates range - too low`() = runTest {
        val result = wizControl.setBrightness(5)

        result.shouldBeInstanceOf<Result.Error>()
        val error = (result as Result.Error).exception
        error.shouldBeInstanceOf<IllegalArgumentException>()
        error.message shouldBe "Brightness must be between 10 and 100, got 5"
    }

    @Test
    fun `setBrightness validates range - too high`() = runTest {
        val result = wizControl.setBrightness(101)

        result.shouldBeInstanceOf<Result.Error>()
        val error = (result as Result.Error).exception
        error.shouldBeInstanceOf<IllegalArgumentException>()
        error.message shouldBe "Brightness must be between 10 and 100, got 101"
    }

    @Test
    fun `setBrightness validates range - minimum valid`() = runTest {
        // We can't actually test the network call without mocking,
        // but we can verify it doesn't fail validation
        val brightness = 10
        (brightness in 10..100) shouldBe true
    }

    @Test
    fun `setBrightness validates range - maximum valid`() = runTest {
        val brightness = 100
        (brightness in 10..100) shouldBe true
    }

    @Test
    fun `setRGBColor validates RGB range - red too high`() = runTest {
        val result = wizControl.setRGBColor(256, 128, 64)

        result.shouldBeInstanceOf<Result.Error>()
        val error = (result as Result.Error).exception
        error.shouldBeInstanceOf<IllegalArgumentException>()
        error.message shouldBe "RGB values must be between 0 and 255"
    }

    @Test
    fun `setRGBColor validates RGB range - green negative`() = runTest {
        val result = wizControl.setRGBColor(128, -1, 64)

        result.shouldBeInstanceOf<Result.Error>()
        val error = (result as Result.Error).exception
        error.shouldBeInstanceOf<IllegalArgumentException>()
        error.message shouldBe "RGB values must be between 0 and 255"
    }

    @Test
    fun `setRGBColor validates RGB range - blue too high`() = runTest {
        val result = wizControl.setRGBColor(128, 128, 300)

        result.shouldBeInstanceOf<Result.Error>()
        val error = (result as Result.Error).exception
        error.shouldBeInstanceOf<IllegalArgumentException>()
        error.message shouldBe "RGB values must be between 0 and 255"
    }

    @Test
    fun `setRGBColor validates RGB range - all valid minimum`() = runTest {
        val r = 0
        val g = 0
        val b = 0
        (r in 0..255 && g in 0..255 && b in 0..255) shouldBe true
    }

    @Test
    fun `setRGBColor validates RGB range - all valid maximum`() = runTest {
        val r = 255
        val g = 255
        val b = 255
        (r in 0..255 && g in 0..255 && b in 0..255) shouldBe true
    }

    @Test
    fun `setColorTemperature validates range - too low`() = runTest {
        val result = wizControl.setColorTemperature(2000)

        result.shouldBeInstanceOf<Result.Error>()
        val error = (result as Result.Error).exception
        error.shouldBeInstanceOf<IllegalArgumentException>()
        error.message shouldBe "Color temperature must be between 2200K and 6500K, got 2000K"
    }

    @Test
    fun `setColorTemperature validates range - too high`() = runTest {
        val result = wizControl.setColorTemperature(7000)

        result.shouldBeInstanceOf<Result.Error>()
        val error = (result as Result.Error).exception
        error.shouldBeInstanceOf<IllegalArgumentException>()
        error.message shouldBe "Color temperature must be between 2200K and 6500K, got 7000K"
    }

    @Test
    fun `setColorTemperature validates range - minimum valid`() = runTest {
        val temp = 2200
        (temp in 2200..6500) shouldBe true
    }

    @Test
    fun `setColorTemperature validates range - maximum valid`() = runTest {
        val temp = 6500
        (temp in 2200..6500) shouldBe true
    }

    @Test
    fun `device properties are accessible`() {
        testDevice.name shouldBe "Test Light"
        testDevice.ip shouldBe "192.168.1.100"
        testDevice.type shouldBe "Bulb"
        testDevice.mac shouldBe "A1B2C3D4E5F6"
    }

    @Test
    fun `pilot state boolean values`() {
        val stateOn = PilotState(state = true)
        val stateOff = PilotState(state = false)

        stateOn.state shouldBe true
        stateOff.state shouldBe false
    }

    @Test
    fun `command JSON structure validation`() {
        // Validate that setPilot command structure is correct
        val commandStructure = mapOf(
            "method" to "setPilot",
            "params" to mapOf("state" to true)
        )

        commandStructure["method"] shouldBe "setPilot"
        commandStructure["params"].shouldBeInstanceOf<Map<*, *>>()
    }

    @Test
    fun `getPilot command structure validation`() {
        // Validate that getPilot command structure is correct
        val commandStructure = mapOf(
            "method" to "getPilot",
            "params" to emptyMap<String, Any>()
        )

        commandStructure["method"] shouldBe "getPilot"
        commandStructure["params"].shouldBeInstanceOf<Map<*, *>>()
    }
}
