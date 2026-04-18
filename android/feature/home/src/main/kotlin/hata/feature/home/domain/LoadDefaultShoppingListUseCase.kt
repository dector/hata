package hata.feature.home.domain

import hata.feature.home.repository.DeviceRepository
import javax.inject.Inject


class LoadDefaultShoppingListUseCase @Inject constructor(
    private val deviceRepository: DeviceRepository,
) {

    suspend fun run(): Result<List<ShoppingItem>> {
        return deviceRepository.getDefaultShoppingList()
    }
}
