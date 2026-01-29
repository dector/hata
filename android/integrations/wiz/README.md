# WiZ Light Control - Kotlin/Android Implementation

## Overview

This package provides functionality to discover and control WiZ smart lights on a local network using UDP communication. This is a Kotlin reimplementation of the Go package from `cave/pkg/int/wiz`.

## Features

- Discover WiZ devices on the network
- Turn lights on and off
- Toggle the power state of lights
- Get the current state of a light (power, brightness, color, etc.)
- Set brightness
- Set RGB color
- Set color temperature

## Required Permissions

Add these permissions to your `AndroidManifest.xml`:

```xml
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_WIFI_STATE" />
<uses-permission android:name="android.permission.CHANGE_WIFI_MULTICAST_STATE" />
```

## Usage

### Discovering Devices

To discover devices, use the `WizDiscovery.scan()` function. It listens for responses for a specified duration (default 10 seconds).

```kotlin
import hata.integrations.wiz.WizDiscovery
import kotlinx.coroutines.runBlocking
import kotlin.time.Duration.Companion.seconds

fun main() = runBlocking {
    println("Scanning for WiZ devices for 10 seconds...")
    val devices = WizDiscovery.scan(10.seconds)

    if (devices.isEmpty()) {
        println("No devices found.")
        return@runBlocking
    }

    println("Found ${devices.size} devices:")
    devices.forEach { device ->
        println("  - Name: ${device.name}, IP: ${device.ip}, MAC: ${device.mac}")
    }
}
```

### Controlling Devices

Once you have a `WizDevice` object from the discovery scan, create a `WizControl` instance to control it.

#### Creating a Device Manually

If you already know the IP address of your device, you can create a `WizDevice` instance manually:

```kotlin
import hata.integrations.wiz.WizDevice
import hata.integrations.wiz.WizControl

val light = WizDevice(
    name = "My Desk Lamp",
    ip = "192.168.1.123",
    type = "Bulb",
    mac = "A1B2C3D4E5F6"
)
val control = WizControl(light)
```

#### Basic Controls

Turn a light on and off:

```kotlin
import hata.integrations.wiz.Result

// Turn on
when (val result = control.turnOn()) {
    is Result.Success -> println("Light turned on")
    is Result.Error -> println("Failed to turn on: ${result.exception.message}")
}

// Turn off
when (val result = control.turnOff()) {
    is Result.Success -> println("Light turned off")
    is Result.Error -> println("Failed to turn off: ${result.exception.message}")
}

// Toggle
control.toggle()
```

#### Setting Brightness

Brightness can be set on a scale from 10 to 100:

```kotlin
// Set brightness to 50%
when (val result = control.setBrightness(50)) {
    is Result.Success -> println("Brightness set to 50%")
    is Result.Error -> println("Failed: ${result.exception.message}")
}
```

#### Setting Color

You can set the color using either RGB values or color temperature.

**RGB Color:**
Values for Red, Green, and Blue range from 0 to 255.

```kotlin
// Set color to a nice purple
when (val result = control.setRGBColor(128, 0, 255)) {
    is Result.Success -> println("Color set to purple")
    is Result.Error -> println("Failed: ${result.exception.message}")
}
```

**Color Temperature:**
Temperature is measured in Kelvin, from 2200K (warm white) to 6500K (cool white).

```kotlin
// Set a warm white color
when (val result = control.setColorTemperature(2700)) {
    is Result.Success -> println("Color temperature set to warm white")
    is Result.Error -> println("Failed: ${result.exception.message}")
}
```

#### Getting Device State

You can retrieve the current state of the light:

```kotlin
when (val result = control.getState()) {
    is Result.Success -> {
        val state = result.data
        println("Light State:")
        println("  - On: ${state.state}")
        println("  - Brightness: ${state.dimming}")
        println("  - Color Temp: ${state.temp}K")
        println("  - RGB: (${state.r}, ${state.g}, ${state.b})")
    }
    is Result.Error -> println("Failed to get state: ${result.exception.message}")
}
```

## Architecture

- **WizModels.kt**: Data classes for devices and communication
- **WizDiscovery.kt**: Device discovery via UDP broadcast
- **WizControl.kt**: Device control commands
- **Result.kt**: Result wrapper for async operations

## Implementation Details

- Uses Kotlin Coroutines for async operations
- UDP communication on port 38899
- JSON serialization with kotlinx.serialization
- No external HTTP libraries needed (uses DatagramSocket directly)

## Differences from Go Implementation

1. Uses Kotlin coroutines instead of Go goroutines
2. Result type for error handling instead of Go's error return values
3. Uses kotlinx.serialization instead of encoding/json
4. Android-specific logging with Log instead of Go's log package
5. No OkHttp dependency needed for UDP communication

## Module Structure

This integration is packaged as a separate Gradle module:

```
integrations/wiz/
├── build.gradle.kts
├── src/
│   ├── main/
│   │   ├── kotlin/hata/integrations/wiz/
│   │   │   ├── WizModels.kt
│   │   │   ├── WizDiscovery.kt
│   │   │   ├── WizControl.kt
│   │   │   └── Result.kt
│   │   └── AndroidManifest.xml
│   └── test/
│       └── kotlin/hata/integrations/wiz/
│           ├── WizModelsTest.kt
│           ├── WizDiscoveryTest.kt
│           └── WizControlTest.kt
```

## Adding to Your Project

Add the module dependency in your app's `build.gradle.kts`:

```kotlin
dependencies {
    implementation(project(":integrations:wiz"))
}
```
