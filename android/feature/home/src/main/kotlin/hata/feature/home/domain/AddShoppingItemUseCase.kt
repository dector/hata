package hata.feature.home.domain

import hata.feature.home.repository.DeviceRepository
import javax.inject.Inject

class AddShoppingItemUseCase @Inject constructor(
    private val deviceRepository: DeviceRepository,
) {

    suspend fun run(name: String): Result<ShoppingItem> {
        return deviceRepository.addDefaultShoppingItem(name = name)
    }
}
