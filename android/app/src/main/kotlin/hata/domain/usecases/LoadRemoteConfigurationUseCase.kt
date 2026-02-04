package hata.domain.usecases

import hata.data.api.RemoteConfigurationApi
import hata.data.models.Device
import hata.data.models.DeviceType
import hata.ui.components.DeviceCard
import hata.ui.home.HomeDisplay
import hata.ui.home.HomeDisplayData
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext


class LoadRemoteConfigurationUseCase(
    private val configurationApi: RemoteConfigurationApi,
    private val dispatcher: CoroutineDispatcher = Dispatchers.IO,
) {

    suspend operator fun invoke(): Result<HomeDisplayData> = withContext(dispatcher) {
        configurationApi
            .loadConfiguration()
            .mapCatching { configuration ->
                val home = configuration.homes.firstOrNull()
                val devices = configuration.homes
                    .flatMap { it.rooms }
                    .flatMap { room ->
                        room.devices.map { device -> room to device }
                    }
                    .map { (room, device) ->
                        DeviceCard.Generic(
                            title = "[${room.name}] ${device.name}",
                            status = toStatus(device),
                            isOn = false,
                            icon = when (device.type) {
                                DeviceType.Light -> DeviceCard.Icon.Light
                                DeviceType.Unknown -> DeviceCard.Icon.Default
                            },
                        )
                    }

                HomeDisplayData(
                    home = HomeDisplay(name = home?.name ?: "Home"),
                    devices = devices,
                )
            }
    }
}

private fun toStatus(device: Device): String {
    val integration = device.integration
        ?: return "Ready"

    val integrationName = integration.id
        .substringBefore('-')
        .replaceFirstChar { it.uppercase() }
    return "[$integrationName] ${integration.data.ip}"
}
