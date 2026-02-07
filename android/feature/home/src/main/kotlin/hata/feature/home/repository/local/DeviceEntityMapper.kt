package hata.feature.home.repository.local

import hata.data.models.Device
import hata.data.models.DeviceState


internal fun DeviceEntity.toDomain(): Device {
    val normalizedState = runCatching { DeviceState.valueOf(state) }
        .getOrDefault(DeviceState.Unknown)

    return Device(
        id = id,
        name = name,
        integrationId = integrationId,
        integrationData = IntegrationDataSerializer.deserialize(integrationDataJson),
        state = normalizedState,
        houseId = houseId,
    )
}

internal fun Device.toEntity(): DeviceEntity {
    return DeviceEntity(
        id = id,
        name = name,
        integrationId = integrationId,
        integrationDataJson = IntegrationDataSerializer.serialize(integrationData),
        state = state.name,
        houseId = houseId,
    )
}
