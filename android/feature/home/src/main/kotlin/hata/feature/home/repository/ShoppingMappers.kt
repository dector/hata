package hata.feature.home.repository

import hata.data.api.models.ApiShoppingItemInfo
import hata.feature.home.domain.ShoppingItem


internal fun ApiShoppingItemInfo.toShoppingItem(): ShoppingItem {
    return ShoppingItem(
        id = uid,
        name = name,
        isChecked = checkedAt != null,
    )
}
