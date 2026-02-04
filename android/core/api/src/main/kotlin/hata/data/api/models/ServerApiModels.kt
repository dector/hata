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
