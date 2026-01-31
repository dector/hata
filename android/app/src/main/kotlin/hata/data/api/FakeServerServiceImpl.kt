package hata.data.api

import hata.data.models.ServerInfo
import kotlinx.coroutines.delay


internal class FakeServerServiceImpl : ServerService {

    override suspend fun connect(serverUrl: String): Result<ServerInfo> {
        // Simulate network delay
        delay(1000)

        return if (serverUrl == "http://localhost") {
            Result.success(
                ServerInfo(
                    serverUrl = serverUrl,
                    serverName = "Local Dev Server",
                    version = "1.0.0-dev",
                )
            )
        } else {
            Result.failure(
                Exception("Failed to connect. Only 'http://localhost' allowed in debug mode.")
            )
        }
    }

    override suspend fun login(
        serverUrl: String,
        username: String,
        password: String,
    ): Result<String> {
        // Simulate network delay
        delay(1000)

        return if (username == "admin" && password == "admin") {
            Result.success("fake_auth_token_${System.currentTimeMillis()}")
        } else {
            Result.failure(Exception("Invalid credentials. Use admin/admin"))
        }
    }
}
