package hata.data.models


data class Session(
    val token: String,
    val user: SessionUser,
    val isActive: Boolean,
)

data class SessionUser(
    val name: String,
)
