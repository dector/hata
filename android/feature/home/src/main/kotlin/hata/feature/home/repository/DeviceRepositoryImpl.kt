package hata.feature.home.repository

import hata.core.annotations.IoDispatcher
import hata.data.api.ServerApi
import hata.data.api.ServerServiceFactory
import hata.data.api.models.ApiDevice
import hata.data.models.Device
import hata.data.models.DeviceState
import hata.feature.home.domain.NewState
import hata.feature.home.domain.ShoppingItem
import hata.feature.home.repository.local.toDomain
import hata.feature.home.repository.local.toEntity
import hata.feature.home.repository.local.UserDevicesDatabaseProvider
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
    private val userDevicesDatabaseProvider: UserDevicesDatabaseProvider,
    @param:IoDispatcher
    private val dispatcher: CoroutineDispatcher,
) : DeviceRepository {

    private var serverApi: ServerApi? = null

    override suspend fun listByHouse(houseId: String): List<Device> = withContext(dispatcher) {
        val deviceDao = userDevicesDatabaseProvider.deviceDao()
        deviceDao
            .listByHouse(houseId)
            .map { it.toDomain() }
    }

    override suspend fun listByUser(): List<Device> = withContext(dispatcher) {
        val deviceDao = userDevicesDatabaseProvider.deviceDao()
        deviceDao
            .listByUser()
            .map { it.toDomain() }
    }

    override suspend fun syncByUser(): Result<Unit> = withContext(dispatcher) {
        runCatching {
            val session = requireSession()
            val service = getServerService(session.serverUrl)
            val apiDevices = service.fetchDevicesByUser().getOrThrow()
            val devices = updateLocalStates(apiDevices.map { it.toDevice() })

            val deviceDao = userDevicesDatabaseProvider.deviceDao()
            deviceDao.replaceAll(devices.map { it.toEntity() })
        }
    }

    override suspend fun getDefaultShoppingList(): Result<List<ShoppingItem>> = withContext(dispatcher) {
        runCatching {
            val session = requireSession()
            val service = getServerService(session.serverUrl)
            service.fetchDefaultShoppingListItems()
                .getOrThrow()
                .map { it.toShoppingItem() }
        }
    }

    override suspend fun setShoppingItemChecked(itemId: String, checked: Boolean): Result<ShoppingItem> = withContext(dispatcher) {
        runCatching {
            val session = requireSession()
            val service = getServerService(session.serverUrl)
            service.setDefaultShoppingItemChecked(itemId = itemId, checked = checked)
                .getOrThrow()
                .toShoppingItem()
        }
    }

    override suspend fun toggleDevice(deviceId: String, newState: NewState): Result<Device> =
        withContext(dispatcher) {
            runCatching {
                val session = requireSession()
                val service = getServerService(session.serverUrl)
                val deviceDao = userDevicesDatabaseProvider.deviceDao()
                val currentDevice = deviceDao.findById(deviceId)?.toDomain()
                    ?: throw IllegalArgumentException("Device with id '$deviceId' is not available")
                val houseId = currentDevice.houseId
                    ?: throw IllegalStateException("Device with id '$deviceId' has no house ID")

                val requestedState = when (newState) {
                    NewState.On -> "on"
                    NewState.Off -> "off"
                }

                val serverState = service.setDeviceState(
                    houseId = houseId,
                    deviceId = deviceId,
                    newState = requestedState,
                ).getOrThrow()

                val updatedDevice = currentDevice.copy(
                    state = serverState.toDeviceState(),
                )
                deviceDao.upsert(updatedDevice.toEntity())
                updatedDevice
            }
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
                    state = if (stateResult.data.state) DeviceState.On else DeviceState.Off,
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
    return Device(
        id = id,
        name = name,
        integrationId = integration?.id ?: "unknown",
        integrationData = integration?.data,
        state = state.toDeviceState(),
        houseId = house?.id,
    )
}

private fun String.toDeviceState(): DeviceState {
    return when (lowercase()) {
        "on" -> DeviceState.On
        "off" -> DeviceState.Off
        else -> DeviceState.Unknown
    }
}
