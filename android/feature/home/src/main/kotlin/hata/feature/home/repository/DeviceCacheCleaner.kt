package hata.feature.home.repository


interface DeviceCacheCleaner {

    suspend fun clearCurrentUserCache()
}
