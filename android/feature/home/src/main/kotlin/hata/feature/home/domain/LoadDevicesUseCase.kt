package hata.feature.home.domain

import hata.data.models.Device
import hata.feature.home.repository.DeviceRepository
import javax.inject.Inject


class LoadDevicesUseCase @Inject constructor(
    private val deviceRepository: DeviceRepository,
) {

    suspend fun run(): LoadDevicesResult {
        val cachedDevices = deviceRepository.listByUser()
        val syncResult = deviceRepository.syncByUser()
        val latestDevices = deviceRepository.listByUser()

        if (latestDevices.isEmpty() && syncResult.isFailure) {
            throw (syncResult.exceptionOrNull() ?: IllegalStateException("Failed to sync devices"))
        }

        return LoadDevicesResult(
            devices = latestDevices.ifEmpty { cachedDevices },
            syncError = syncResult.exceptionOrNull(),
        )
    }
}

data class LoadDevicesResult(
    val devices: List<Device>,
    val syncError: Throwable?,
)
