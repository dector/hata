package hata.data.api


interface ServerServiceFactory {

    fun create(serverUrl: String): ServerApi
}
