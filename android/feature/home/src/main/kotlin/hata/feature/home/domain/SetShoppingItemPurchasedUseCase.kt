package hata.feature.home.domain

import hata.feature.home.repository.DeviceRepository
import javax.inject.Inject

class SetShoppingItemPurchasedUseCase @Inject constructor(
    private val deviceRepository: DeviceRepository,
) {

    suspend fun run(itemId: String): Result<ShoppingItem> {
        return deviceRepository.setShoppingItemPurchased(itemId = itemId)
    }
}
