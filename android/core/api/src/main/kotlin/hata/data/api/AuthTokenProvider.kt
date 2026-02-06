package hata.data.api


interface AuthTokenProvider {

    fun getToken(): String?
}
