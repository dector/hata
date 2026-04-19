package hata.feature.home.ui

import hata.data.models.Device
import hata.feature.home.domain.AddShoppingItemUseCase
import hata.feature.home.domain.LoadDefaultShoppingListUseCase
import hata.feature.home.domain.NewState
import hata.feature.home.domain.SetShoppingItemPurchasedUseCase
import hata.feature.home.domain.ShoppingItem
import hata.feature.home.repository.DeviceRepository
import hata.navigation.Navigator
import hata.navigation.Route
import io.kotest.matchers.shouldBe
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.async
import kotlinx.coroutines.flow.first
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
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = defaultItems,
        )
    }

    @Test
    fun `init sorts items by unchecked first and name alphabetically`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(
                listOf(
                    ShoppingItem(id = "4", name = "Bread", isChecked = true),
                    ShoppingItem(id = "2", name = "Apple", isChecked = false),
                    ShoppingItem(id = "1", name = "milk", isChecked = false),
                    ShoppingItem(id = "3", name = "Cheese", isChecked = true),
                ),
            ),
        )
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = listOf(
                ShoppingItem(id = "2", name = "Apple", isChecked = false),
                ShoppingItem(id = "1", name = "milk", isChecked = false),
                ShoppingItem(id = "4", name = "Bread", isChecked = true),
                ShoppingItem(id = "3", name = "Cheese", isChecked = true),
            ),
        )
    }

    @Test
    fun `init load failure updates state to error`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.failure(IllegalStateException("Boom")),
        )
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Error("Boom")
    }

    @Test
    fun `set checked success updates item`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(defaultItems),
            setCheckedResult = Result.success(
                ShoppingItem(id = "1", name = "Milk", isChecked = true),
            ),
        )
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()

        val event = async { vm.events.first() }
        vm.onDispatch(ShoppingUiAction.SetChecked(itemId = "1", checked = true))
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = listOf(
                ShoppingItem(id = "1", name = "Milk", isChecked = true),
            ),
            errorMessage = null,
        )
        event.await() shouldBe ShoppingUiEvent.ItemCheckedChanged(
            itemId = "1",
            itemName = "Milk",
            checked = true,
        )
    }

    @Test
    fun `set unchecked success updates item`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(
                listOf(ShoppingItem(id = "1", name = "Milk", isChecked = true)),
            ),
            setCheckedResult = Result.success(
                ShoppingItem(id = "1", name = "Milk", isChecked = false),
            ),
        )
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()
        vm.onDispatch(ShoppingUiAction.SetChecked(itemId = "1", checked = false))
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = listOf(
                ShoppingItem(id = "1", name = "Milk", isChecked = false),
            ),
            errorMessage = null,
        )
    }

    @Test
    fun `set checked success resorts list`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(
                listOf(
                    ShoppingItem(id = "1", name = "Apple", isChecked = false),
                    ShoppingItem(id = "2", name = "Banana", isChecked = false),
                ),
            ),
            setCheckedResult = Result.success(
                ShoppingItem(id = "1", name = "Apple", isChecked = true),
            ),
        )
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()
        vm.onDispatch(ShoppingUiAction.SetChecked(itemId = "1", checked = true))
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = listOf(
                ShoppingItem(id = "2", name = "Banana", isChecked = false),
                ShoppingItem(id = "1", name = "Apple", isChecked = true),
            ),
            errorMessage = null,
        )
    }

    @Test
    fun `set checked failure keeps item and shows error`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(defaultItems),
            setCheckedResult = Result.failure(IllegalStateException("No internet")),
        )
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()
        vm.onDispatch(ShoppingUiAction.SetChecked(itemId = "1", checked = true))
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = defaultItems,
            errorMessage = "No internet",
        )
    }

    @Test
    fun `save new item success appends and resorts then closes dialog`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(
                listOf(
                    ShoppingItem(id = "1", name = "Milk", isChecked = false),
                ),
            ),
            addItemResult = Result.success(
                ShoppingItem(id = "2", name = "Apple", isChecked = false),
            ),
        )
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()
        vm.onDispatch(ShoppingUiAction.OpenAddDialog)
        vm.onDispatch(ShoppingUiAction.UpdateNewItemName("Apple"))
        vm.onDispatch(ShoppingUiAction.SaveNewItem)
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = listOf(
                ShoppingItem(id = "2", name = "Apple", isChecked = false),
                ShoppingItem(id = "1", name = "Milk", isChecked = false),
            ),
            errorMessage = null,
            isAddDialogVisible = false,
            newItemName = "",
            isSavingItem = false,
        )
        repository.addItemCalls shouldBe 1
        repository.lastAddedName shouldBe "Apple"
    }

    @Test
    fun `save new item failure keeps list and shows error with dialog open`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(defaultItems),
            addItemResult = Result.failure(IllegalStateException("Create failed")),
        )
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()
        vm.onDispatch(ShoppingUiAction.OpenAddDialog)
        vm.onDispatch(ShoppingUiAction.UpdateNewItemName("Bread"))
        vm.onDispatch(ShoppingUiAction.SaveNewItem)
        advanceUntilIdle()

        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = defaultItems,
            errorMessage = "Create failed",
            isAddDialogVisible = true,
            newItemName = "Bread",
            isSavingItem = false,
        )
        repository.addItemCalls shouldBe 1
    }

    @Test
    fun `save new item with blank name does not call repository`() = runTest {
        val repository = FakeDeviceRepository(
            shoppingListResult = Result.success(defaultItems),
        )
        val vm = createVm(repository)

        vm.onDispatch(ShoppingUiAction.Init)
        advanceUntilIdle()
        vm.onDispatch(ShoppingUiAction.OpenAddDialog)
        vm.onDispatch(ShoppingUiAction.UpdateNewItemName("   "))
        vm.onDispatch(ShoppingUiAction.SaveNewItem)
        advanceUntilIdle()

        repository.addItemCalls shouldBe 0
        vm.uiState.value shouldBe ShoppingUiState.Loaded(
            items = defaultItems,
            isAddDialogVisible = true,
            newItemName = "   ",
        )
    }

    private fun createVm(repository: FakeDeviceRepository): ShoppingViewModel {
        return ShoppingViewModel(
            loadDefaultShoppingListUseCase = LoadDefaultShoppingListUseCase(repository),
            setShoppingItemPurchasedUseCase = SetShoppingItemPurchasedUseCase(repository),
            addShoppingItemUseCase = AddShoppingItemUseCase(repository),
            navigator = FakeNavigator(),
        )
    }
}

private class FakeDeviceRepository(
    private val shoppingListResult: Result<List<ShoppingItem>>,
    private val setCheckedResult: Result<ShoppingItem> = Result.success(
        ShoppingItem(id = "1", name = "Milk", isChecked = true),
    ),
    private val addItemResult: Result<ShoppingItem> = Result.success(
        ShoppingItem(id = "2", name = "Bread", isChecked = false),
    ),
) : DeviceRepository {

    var addItemCalls: Int = 0
        private set

    var lastAddedName: String? = null
        private set

    override suspend fun listByHouse(houseId: String): List<Device> = emptyList()

    override suspend fun listByUser(): List<Device> = emptyList()

    override suspend fun syncByUser(): Result<Unit> = Result.success(Unit)

    override suspend fun getDefaultShoppingList(): Result<List<ShoppingItem>> = shoppingListResult

    override suspend fun setShoppingItemChecked(itemId: String, checked: Boolean): Result<ShoppingItem> = setCheckedResult

    override suspend fun addDefaultShoppingItem(name: String): Result<ShoppingItem> {
        addItemCalls += 1
        lastAddedName = name
        return addItemResult
    }

    override suspend fun toggleDevice(deviceId: String, newState: NewState): Result<Device> {
        throw UnsupportedOperationException("Not used in tests")
    }
}

private class FakeNavigator : Navigator {
    override fun goTo(route: Route) = Unit

    override fun goBack() = Unit
}
