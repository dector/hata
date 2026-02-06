package hata.feature.home.ui

import hata.data.models.Device
import hata.data.models.DeviceState
import hata.ui.components.DeviceCard


internal fun Device.toDeviceCard(showHouseId: Boolean): DeviceCard.Generic {
    val statusParts = buildList {
        if (integrationId.isNotBlank() && integrationId != "unknown") {
            add(integrationId)
        }
        add(stateLabel())
        if (showHouseId) {
            houseId?.let { add(it) }
        }
    }

    return DeviceCard.Generic(
        id = id,
        title = name,
        status = statusParts.joinToString(" | "),
        isOn = state == DeviceState.On,
        canControl = isControllable(),
        icon = iconForIntegration(integrationId),
    )
}

private fun Device.isControllable(): Boolean {
    if (!integrationId.isWizIntegration()) {
        return true
    }

    val type = integrationData?.get("type") as? String
    val ip = integrationData?.get("ip") as? String

    return type.equals("wifi", ignoreCase = true) && !ip.isNullOrBlank()
}

private fun Device.stateLabel(): String {
    return when (state) {
        DeviceState.On -> "On"
        DeviceState.Off -> "Off"
        DeviceState.Unknown -> "Unknown"
    }
}

private fun iconForIntegration(integrationId: String): DeviceCard.Icon {
    return if (integrationId.isWizIntegration()) {
        DeviceCard.Icon.Light
    } else {
        DeviceCard.Icon.Default
    }
}

private fun String.isWizIntegration(): Boolean {
    return startsWith("wiz", ignoreCase = true)
}
