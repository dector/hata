package hata.feature.home.domain

import hata.data.models.Device
import hata.feature.home.repository.DeviceRepository
import javax.inject.Inject


class LoadDevicesUseCase @Inject constructor(
    private val deviceRepository: DeviceRepository,
) {

    suspend fun run(): Result<LoadDevicesResult> {
        val syncResult = deviceRepository.syncByUser()
        val latestDevices = deviceRepository.listByUser()

        if (latestDevices.isEmpty() && syncResult.isFailure) {
            return Result.failure(
                syncResult.exceptionOrNull() ?: IllegalStateException("Failed to sync devices"),
            )
        }

        return Result.success(
            LoadDevicesResult(
                devices = latestDevices,
                syncError = syncResult.exceptionOrNull(),
            ),
        )
    }
}

data class LoadDevicesResult(
    val devices: List<Device>,
    val syncError: Throwable?,
)
