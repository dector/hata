package hata.feature.home.domain

import hata.feature.home.repository.DeviceRepository
import javax.inject.Inject

class SetShoppingItemPurchasedUseCase @Inject constructor(
    private val deviceRepository: DeviceRepository,
) {

    suspend fun run(itemId: String, checked: Boolean): Result<ShoppingItem> {
        return deviceRepository.setShoppingItemChecked(
            itemId = itemId,
            checked = checked,
        )
    }
}
