package hata.feature.home.repository

import hata.core.annotations.IoDispatcher
import hata.data.api.ServerApi
import hata.data.api.ServerServiceFactory
import hata.data.api.models.ApiDevice
import hata.data.models.Device
import hata.data.models.DeviceState
import hata.feature.session.repository.SessionRepository
import hata.integrations.wiz.WizControl
import hata.integrations.wiz.WizDevice
import hata.integrations.wiz.Result as WizResult
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.withContext
import javax.inject.Inject
import javax.inject.Singleton


@Singleton
class DeviceRepositoryImpl @Inject constructor(
    private val serverServiceFactory: ServerServiceFactory,
    private val sessionRepository: SessionRepository,
    @param:IoDispatcher
    private val dispatcher: CoroutineDispatcher,
) : DeviceRepository {

    private val devicesByHouse = mutableMapOf<String, List<Device>>()
    private var devicesByUser: List<Device> = emptyList()
    private var serverApi: ServerApi? = null

    override suspend fun listByHouse(houseId: String): List<Device> = withContext(dispatcher) {
        val session = requireSession()
        val service = getServerService(session.serverUrl)
        val result = service.fetchDevicesByHouse(houseId)

        result.fold(
            onSuccess = { apiDevices ->
                val devices = updateLocalStates(apiDevices.map { it.toDevice() })
                devicesByHouse[houseId] = devices
                devices
            },
            onFailure = { error ->
                devicesByHouse[houseId]?.takeIf { it.isNotEmpty() } ?: throw error
            },
        )
    }

    override suspend fun listByUser(): List<Device> = withContext(dispatcher) {
        val session = requireSession()
        val service = getServerService(session.serverUrl)
        val result = service.fetchDevicesByUser()

        result.fold(
            onSuccess = { apiDevices ->
                val devices = updateLocalStates(apiDevices.map { it.toDevice() })
                devicesByUser = devices
                devices.groupBy { it.houseId }.forEach { (houseId, items) ->
                    if (houseId != null) {
                        devicesByHouse[houseId] = items
                    }
                }
                devices
            },
            onFailure = { error ->
                devicesByUser.takeIf { it.isNotEmpty() } ?: throw error
            },
        )
    }

    private suspend fun requireSession() = sessionRepository.getSession()
        ?: throw IllegalStateException("No active session")

    private fun getServerService(serverUrl: String): ServerApi {
        val currentService = serverApi
        if (currentService != null && currentService.serverUrl == serverUrl) {
            return currentService
        }

        return serverServiceFactory.create(serverUrl).also { serverApi = it }
    }

    private suspend fun updateLocalStates(devices: List<Device>): List<Device> {
        return devices.map { device ->
            if (!device.isWizIntegration()) {
                return@map device
            }

            val ip = device.integrationData?.get("ip") as? String
                ?: return@map device

            val type = device.integrationData?.get("type") as? String ?: "unknown"
            val mac = device.integrationData?.get("mac") as? String ?: device.id

            val wizDevice = WizDevice(
                name = device.name,
                ip = ip,
                type = type,
                mac = mac,
            )

            when (val stateResult = WizControl(wizDevice).getState()) {
                is WizResult.Success -> device.copy(
                    state = if (stateResult.data.state) DeviceState.ON else DeviceState.OFF,
                )
                is WizResult.Error -> device
            }
        }
    }
}

private fun Device.isWizIntegration(): Boolean {
    return integrationId.startsWith("wiz", ignoreCase = true)
}

private fun ApiDevice.toDevice(): Device {
    val normalizedState = when (state.lowercase()) {
        "on" -> DeviceState.ON
        "off" -> DeviceState.OFF
        else -> DeviceState.UNKNOWN
    }

    return Device(
        id = id,
        name = name,
        integrationId = integration?.id ?: "unknown",
        integrationData = integration?.data,
        state = normalizedState,
        houseId = house?.id,
    )
}
