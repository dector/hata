package hata.feature.home.ui

import hata.data.models.Device
import hata.feature.home.domain.LoadDefaultShoppingListUseCase
import hata.feature.home.domain.NewState
import hata.feature.home.domain.ShoppingItem
import hata.feature.home.repository.DeviceRepository
import hata.navigation.Navigator
import hata.navigation.Route
import io.kotest.matchers.shouldBe
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.runTest
import org.junit.Rule
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class ShoppingViewModelTest {

    @get:Rule
    val mainDispatcherRule = MainDispatcherRule()

    @Test
    fun `init loads items and updates state to loaded`() = runTest {
        val vm = ShoppingViewModel(
            loadDefaultShoppingListUseCase = LoadDefaultShoppingListUseCase(
                deviceRepository = FakeDeviceRepository(
                    shoppingListResult = Result.success(
                        listOf(
                            ShoppingItem(id = "1", name = "Milk", isChecked = false),
                        ),
                    ),
                ),
            ),
            navigator = FakeNavigator(),
        )

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = listOf(
                ShoppingItem(id = "1", name = "Milk", isChecked = false),
            ),
        )
    }

    @Test
    fun `init load failure updates state to error`() = runTest {
        val vm = ShoppingViewModel(
            loadDefaultShoppingListUseCase = LoadDefaultShoppingListUseCase(
                deviceRepository = FakeDeviceRepository(
                    shoppingListResult = Result.failure(IllegalStateException("Boom")),
                ),
            ),
            navigator = FakeNavigator(),
        )

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Error("Boom")
    }
}

private class FakeDeviceRepository(
    private val shoppingListResult: Result<List<ShoppingItem>>,
) : DeviceRepository {
    override suspend fun listByHouse(houseId: String): List<Device> = emptyList()

    override suspend fun listByUser(): List<Device> = emptyList()

    override suspend fun syncByUser(): Result<Unit> = Result.success(Unit)

    override suspend fun getDefaultShoppingList(): Result<List<ShoppingItem>> = shoppingListResult

    override suspend fun toggleDevice(deviceId: String, newState: NewState): Result<Device> {
        throw UnsupportedOperationException("Not used in tests")
    }
}

private class FakeNavigator : Navigator {
    override fun goTo(route: Route) = Unit

    override fun goBack() = Unit
}
