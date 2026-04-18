package hata.data.api.models

import com.squareup.moshi.JsonClass


// Ping response
@JsonClass(generateAdapter = true)
data class PingResponse(
    val serverName: String,
    val version: String,
    val apiVersion: String,
)

// Login request
@JsonClass(generateAdapter = true)
data class LoginRequest(
    val username: String,
    val password: String,
)

// Login response
@JsonClass(generateAdapter = true)
data class LoginResponse(
    val session: SessionInfo,
    val user: UserInfo,
)

@JsonClass(generateAdapter = true)
data class SessionInfo(
    val token: String,
    val validUntil: String, // ISO 8601 format
)

@JsonClass(generateAdapter = true)
data class UserInfo(
    val displayName: String,
)

// Error response
@JsonClass(generateAdapter = true)
data class ErrorResponse(
    val error: ErrorDetail,
)

@JsonClass(generateAdapter = true)
data class ErrorDetail(
    val message: String,
    val code: String,
)

// House response
@JsonClass(generateAdapter = true)
data class ApiHouseListResponse(
    val houses: List<ApiHouseInfo> = emptyList(),
)

@JsonClass(generateAdapter = true)
data class ApiHouseInfo(
    val id: String,
    val displayName: String,
    val role: String,
)

// Shopping list responses
@JsonClass(generateAdapter = true)
data class ApiShoppingListListResponse(
    val shoppingLists: List<ApiShoppingListInfo> = emptyList(),
)

@JsonClass(generateAdapter = true)
data class ApiShoppingListInfo(
    val uid: String,
    val name: String,
    val createdAt: String,
    val updatedAt: String,
)

@JsonClass(generateAdapter = true)
data class ApiShoppingItemListResponse(
    val items: List<ApiShoppingItemInfo> = emptyList(),
)

@JsonClass(generateAdapter = true)
data class ApiShoppingItemInfo(
    val uid: String,
    val name: String,
    val position: Int,
    val checkedAt: String? = null,
    val checkedByUserId: Int? = null,
    val deletedAt: String? = null,
    val createdAt: String,
    val updatedAt: String,
)
