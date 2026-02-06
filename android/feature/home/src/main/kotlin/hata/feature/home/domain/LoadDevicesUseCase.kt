package hata.feature.home.domain

import hata.feature.home.repository.DeviceRepository
import javax.inject.Inject


class LoadDevicesUseCase @Inject constructor(
    private val deviceRepository: DeviceRepository,
) {

    suspend fun run() = deviceRepository.listByUser()
}
