package hata.data.api


class RealServerServiceFactory(
    private val tokenProvider: AuthTokenProvider,
) : ServerServiceFactory {

    override fun create(serverUrl: String): ServerApi =
        RealServerApi(
            serverUrl = serverUrl,
            tokenProvider = tokenProvider,
        )
}
