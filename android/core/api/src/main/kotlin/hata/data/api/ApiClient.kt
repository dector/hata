package hata.data.api

import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import hata.data.api.models.ErrorResponse
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.HttpException
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory
import java.util.concurrent.TimeUnit


internal object ApiClient {

    private val moshi = Moshi.Builder()
        .add(KotlinJsonAdapterFactory())
        .build()

    fun createRetrofit(
        baseUrl: String,
        tokenProvider: AuthTokenProvider? = null,
    ): Retrofit {
        val normalizedBaseUrl = if (baseUrl.endsWith("/")) baseUrl else "$baseUrl/"

        return Retrofit.Builder()
            .baseUrl(normalizedBaseUrl)
            .client(createOkHttpClient(tokenProvider))
            .addConverterFactory(MoshiConverterFactory.create(moshi))
            .build()
    }

    fun parseErrorMessage(exception: HttpException): String? {
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

private fun createOkHttpClient(
    tokenProvider: AuthTokenProvider?,
): OkHttpClient {
    val builder = OkHttpClient.Builder()
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)

    if (tokenProvider != null) {
        builder.addInterceptor(AuthInterceptor(tokenProvider))
    }

    builder.setLogging()

    return builder.build()
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
