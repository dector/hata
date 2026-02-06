package hata.feature.home.domain

import hata.data.models.Device
import hata.feature.home.repository.DeviceRepository
import javax.inject.Inject


class ToggleDeviceUseCase @Inject constructor(
    private val deviceRepository: DeviceRepository,
) {

    suspend fun run(deviceId: String, newState: NewState): Result<Device> {
        return deviceRepository.toggleDevice(deviceId, newState)
    }
}
