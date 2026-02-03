package hata.data.api

import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import hata.BuildConfig
import hata.data.api.models.ErrorResponse
import hata.data.api.models.LoginRequest
import hata.data.models.ServerInfo
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.HttpException
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory
import java.io.IOException
import java.util.concurrent.TimeUnit

internal class RealServerServiceImpl : ServerService {

    private val moshi = Moshi.Builder()
        .add(KotlinJsonAdapterFactory())
        .build()

    private val okHttpClient = OkHttpClient.Builder()
        .setLogging()
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)
        .build()

    // Emulator: http://10.0.2.2:8080
    override suspend fun connect(serverUrl: String): Result<ServerInfo> {
        return try {
            val retrofit = createRetrofit(serverUrl)
            val api = retrofit.create(HataApi::class.java)

            val response = api.ping()

            Result.success(
                ServerInfo(
                    serverUrl = serverUrl,
                    serverName = response.serverName,
                    version = response.version,
                ),
            )
        } catch (e: HttpException) {
            Result.failure(
                Exception(parseErrorMessage(e) ?: "Failed to connect to server"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Connection failed: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun login(
        serverUrl: String,
        username: String,
        password: String,
    ): Result<String> {
        return try {
            val retrofit = createRetrofit(serverUrl)
            val api = retrofit.create(HataApi::class.java)

            val response = api.login(
                LoginRequest(
                    username = username,
                    password = password,
                ),
            )

            Result.success(response.session.token)
        } catch (e: HttpException) {
            Result.failure(
                Exception(parseErrorMessage(e) ?: "Login failed"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Login failed: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchLatestHouse(
        serverUrl: String,
        token: String,
    ): Result<Unit> {
        return try {
            val retrofit = createRetrofit(serverUrl)
            val api = retrofit.create(HataApi::class.java)

            api.latestHouse("Bearer $token")

            Result.success(Unit)
        } catch (e: HttpException) {
            Result.failure(
                Exception(parseErrorMessage(e) ?: "Failed to load latest house"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load latest house: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    private fun createRetrofit(baseUrl: String): Retrofit {
        // Ensure baseUrl ends with /
        val normalizedBaseUrl = if (baseUrl.endsWith("/")) baseUrl else "$baseUrl/"

        return Retrofit.Builder()
            .baseUrl(normalizedBaseUrl)
            .client(okHttpClient)
            .addConverterFactory(MoshiConverterFactory.create(moshi))
            .build()
    }

    private fun parseErrorMessage(exception: HttpException): String? {
        return try {
            val errorBody = exception.response()?.errorBody()?.string()
            if (errorBody != null) {
                val adapter = moshi.adapter(ErrorResponse::class.java)
                val errorResponse = adapter.fromJson(errorBody)
                errorResponse?.error?.message
            } else {
                null
            }
        } catch (e: Exception) {
            null
        }
    }
}

private fun OkHttpClient.Builder.setLogging(
    applyIf: Boolean = BuildConfig.DEBUG,
): OkHttpClient.Builder = apply {
    if (!applyIf) return@apply

    addInterceptor(
        HttpLoggingInterceptor().apply {
            level = HttpLoggingInterceptor.Level.BODY
        },
    )
}
