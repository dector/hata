package hata.feature.home.ui

import hata.data.models.Device
import hata.feature.home.domain.LoadDefaultShoppingListUseCase
import hata.feature.home.domain.NewState
import hata.feature.home.domain.SetShoppingItemPurchasedUseCase
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

    private val defaultItems = listOf(
        ShoppingItem(id = "1", name = "Milk", isChecked = false),
    )

    @Test
    fun `init loads items and updates state to loaded`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(defaultItems),
        )
        val vm = ShoppingViewModel(
            loadDefaultShoppingListUseCase = LoadDefaultShoppingListUseCase(
                deviceRepository = repository,
            ),
            setShoppingItemPurchasedUseCase = SetShoppingItemPurchasedUseCase(
                deviceRepository = repository,
            ),
            navigator = FakeNavigator(),
        )

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = defaultItems,
        )
    }

    @Test
    fun `init load failure updates state to error`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.failure(IllegalStateException("Boom")),
        )
        val vm = ShoppingViewModel(
            loadDefaultShoppingListUseCase = LoadDefaultShoppingListUseCase(
                deviceRepository = repository,
            ),
            setShoppingItemPurchasedUseCase = SetShoppingItemPurchasedUseCase(
                deviceRepository = repository,
            ),
            navigator = FakeNavigator(),
        )

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Error("Boom")
    }

    @Test
    fun `mark purchased success updates item`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(defaultItems),
            setPurchasedResult = Result.success(
                ShoppingItem(id = "1", name = "Milk", isChecked = true),
            ),
        )
        val vm = ShoppingViewModel(
            loadDefaultShoppingListUseCase = LoadDefaultShoppingListUseCase(repository),
            setShoppingItemPurchasedUseCase = SetShoppingItemPurchasedUseCase(repository),
            navigator = FakeNavigator(),
        )

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()
        vm.onDispatch(ShoppingUiAction.MarkPurchased("1"))
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = listOf(
                ShoppingItem(id = "1", name = "Milk", isChecked = true),
            ),
            errorMessage = null,
        )
    }

    @Test
    fun `mark purchased failure keeps item and shows error`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(defaultItems),
            setPurchasedResult = Result.failure(IllegalStateException("No internet")),
        )
        val vm = ShoppingViewModel(
            loadDefaultShoppingListUseCase = LoadDefaultShoppingListUseCase(repository),
            setShoppingItemPurchasedUseCase = SetShoppingItemPurchasedUseCase(repository),
            navigator = FakeNavigator(),
        )

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()
        vm.onDispatch(ShoppingUiAction.MarkPurchased("1"))
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = defaultItems,
            errorMessage = "No internet",
        )
    }
}

private class FakeDeviceRepository(
    private val shoppingListResult: Result<List<ShoppingItem>>,
    private val setPurchasedResult: Result<ShoppingItem> = Result.success(
        ShoppingItem(id = "1", name = "Milk", isChecked = true),
    ),
) : DeviceRepository {
    override suspend fun listByHouse(houseId: String): List<Device> = emptyList()

    override suspend fun listByUser(): List<Device> = emptyList()

    override suspend fun syncByUser(): Result<Unit> = Result.success(Unit)

    override suspend fun getDefaultShoppingList(): Result<List<ShoppingItem>> = shoppingListResult

    override suspend fun setShoppingItemPurchased(itemId: String): Result<ShoppingItem> = setPurchasedResult

    override suspend fun toggleDevice(deviceId: String, newState: NewState): Result<Device> {
        throw UnsupportedOperationException("Not used in tests")
    }
}

private class FakeNavigator : Navigator {
    override fun goTo(route: Route) = Unit

    override fun goBack() = Unit
}
