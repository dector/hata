package hata.feature.home.domain

import hata.data.models.Device
import hata.feature.home.repository.DeviceRepository
import javax.inject.Inject


class LoadCachedDevicesUseCase @Inject constructor(
    private val deviceRepository: DeviceRepository,
) {

    suspend fun run(): List<Device> {
        return deviceRepository.listByUser()
    }
}
