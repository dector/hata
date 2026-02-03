package hata.feature.session.model


data class Session(
    val token: String,
    val user: SessionUser,
    val isActive: Boolean,
    val serverUrl: String,
    val createdAt: Long,
)

data class SessionUser(
    val name: String,
)
