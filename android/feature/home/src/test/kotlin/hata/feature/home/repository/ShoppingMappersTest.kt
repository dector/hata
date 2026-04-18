package hata.feature.home.repository

import hata.data.api.models.ApiShoppingItemInfo
import io.kotest.matchers.shouldBe
import org.junit.Test


class ShoppingMappersTest {

    @Test
    fun `toShoppingItem maps checked state from checkedAt`() {
        val checkedApiItem = ApiShoppingItemInfo(
            uid = "item-1",
            name = "Milk",
            position = 1,
            checkedAt = "2026-01-01T10:00:00Z",
            checkedByUserId = 1,
            deletedAt = null,
            createdAt = "2026-01-01T09:00:00Z",
            updatedAt = "2026-01-01T10:00:00Z",
        )
        val uncheckedApiItem = checkedApiItem.copy(
            uid = "item-2",
            checkedAt = null,
            checkedByUserId = null,
        )

        checkedApiItem.toShoppingItem().isChecked shouldBe true
        uncheckedApiItem.toShoppingItem().isChecked shouldBe false
    }
}
