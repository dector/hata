package hata.feature.login.domain

import hata.data.api.ServerServiceFactory
import javax.inject.Inject


class LoginUseCase @Inject constructor(
    private val serverServiceFactory: ServerServiceFactory,
) {

    suspend fun run(
        serverUrl: String,
        username: String,
        password: String,
    ): Result<String> {
        val service = serverServiceFactory.create(serverUrl)
        return service
            .login(username, password)
            .onSuccess { service.fetchHouse() }
    }
}
