package hata.data.api

import hata.data.api.models.ApiDevice
import hata.data.api.models.ApiDeviceSetStateRequest
import hata.data.api.models.ApiHouseInfo
import hata.data.api.models.ApiShoppingItemCheckRequest
import hata.data.api.models.ApiShoppingItemCreateRequest
import hata.data.api.models.ApiShoppingItemInfo
import hata.data.api.models.ApiShoppingListInfo
import hata.data.api.models.LoginRequest
import retrofit2.HttpException
import java.io.IOException


class RealServerApi(
    override val serverUrl: String,
    tokenProvider: AuthTokenProvider,
) : ServerApi {

    private val api = ApiClient
        .createRetrofit(serverUrl, tokenProvider)
        .create(HataApi::class.java)

    override suspend fun login(
        username: String,
        password: String,
    ): Result<String> {
        return try {
            val response = api.login(
                LoginRequest(
                    username = username,
                    password = password,
                ),
            )

            Result.success(response.session.token)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Login failed"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Login failed: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchHouse(): Result<Unit> {
        return try {
            api.latestHouse()

            Result.success(Unit)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load latest house"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load latest house: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchHouses(): Result<List<ApiHouseInfo>> {
        return try {
            val response = api.latestHouse()

            Result.success(response.houses)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load houses"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load houses: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchDevicesByHouse(
        houseId: String,
    ): Result<List<ApiDevice>> {
        return try {
            val response = api.devicesByHouse(
                houseId = houseId,
            )

            Result.success(response.devices)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load devices for house"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load devices: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchDevicesByUser(
    ): Result<List<ApiDevice>> {
        return try {
            val response = api.devicesByUser()

            Result.success(response.devices)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load devices"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load devices: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun setDeviceState(
        houseId: String,
        deviceId: String,
        newState: String,
    ): Result<String> {
        return try {
            val response = api.setDeviceState(
                houseId = houseId,
                deviceId = deviceId,
                request = ApiDeviceSetStateRequest(state = newState),
            )

            Result.success(response.state)
        } catch (e: HttpException) {
            val apiError = ApiClient.parseError(e)
            val message = apiError?.message ?: "Failed to set device state"

            val mappedError = when (apiError?.code) {
                DeviceNoAckException.CODE -> DeviceNoAckException(message)
                null -> Exception(message)
                else -> ApiException(message = message, code = apiError.code)
            }

            Result.failure(mappedError)
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to set device state: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchShoppingListsByHouse(
        houseId: String,
    ): Result<List<ApiShoppingListInfo>> {
        return try {
            val response = api.shoppingListsByHouse(houseId = houseId)

            Result.success(response.shoppingLists)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load shopping lists"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load shopping lists: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchShoppingItemsByList(
        houseId: String,
        listId: String,
    ): Result<List<ApiShoppingItemInfo>> {
        return try {
            val response = api.shoppingItemsByList(
                houseId = houseId,
                listId = listId,
            )

            Result.success(response.items)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load shopping items"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load shopping items: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchDefaultShoppingListItems(): Result<List<ApiShoppingItemInfo>> {
        return try {
            val (houseId, listId) = resolveDefaultShoppingListIds()
                ?: return Result.success(emptyList())

            val items = api.shoppingItemsByList(
                houseId = houseId,
                listId = listId,
            ).items

            Result.success(items)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load default shopping list"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load default shopping list: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun setDefaultShoppingItemChecked(
        itemId: String,
        checked: Boolean,
    ): Result<ApiShoppingItemInfo> {
        return try {
            val (houseId, listId) = resolveDefaultShoppingListIds()
                ?: return Result.failure(Exception("Default shopping list is not available"))

            val item = api.setShoppingItemChecked(
                houseId = houseId,
                listId = listId,
                itemId = itemId,
                request = ApiShoppingItemCheckRequest(checked = checked),
            ).item

            Result.success(item)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to update shopping item"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to update shopping item: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun addDefaultShoppingItem(name: String): Result<ApiShoppingItemInfo> {
        return try {
            val (houseId, listId) = resolveDefaultShoppingListIds()
                ?: return Result.failure(Exception("Default shopping list is not available"))

            val item = api.createShoppingItem(
                houseId = houseId,
                listId = listId,
                request = ApiShoppingItemCreateRequest(name = name),
            ).item

            Result.success(item)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to add shopping item"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to add shopping item: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    private suspend fun resolveDefaultShoppingListIds(): Pair<String, String>? {
        val houses = api.latestHouse().houses
        val houseId = houses.firstOrNull()?.id ?: return null

        val lists = api.shoppingListsByHouse(houseId = houseId).shoppingLists
        val listId = lists.firstOrNull { it.uid == DEFAULT_SHOPPING_LIST_UID }?.uid
            ?: lists.firstOrNull()?.uid
            ?: return null

        return houseId to listId
    }

    private companion object {
        const val DEFAULT_SHOPPING_LIST_UID = "default"
    }
}

