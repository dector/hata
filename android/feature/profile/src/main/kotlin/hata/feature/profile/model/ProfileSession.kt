package hata.feature.profile.model


data class ProfileSession(
    val userName: String,
    val isActive: Boolean,
    val serverUrl: String,
    val createdAt: Long,
)
