package hata.data.api

import okhttp3.Interceptor
import okhttp3.Response


class AuthInterceptor(
    private val tokenProvider: AuthTokenProvider,
) : Interceptor {

    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        if (request.header(AUTH_HEADER) != null) {
            return chain.proceed(request)
        }

        val token = tokenProvider
            .getToken()
            ?.takeIf(String::isNotBlank)
            ?: return chain.proceed(request)

        val updatedRequest = request
            .newBuilder()
            .addHeader(AUTH_HEADER, "Bearer $token")
            .build()

        return chain.proceed(updatedRequest)
    }

    private companion object {
        const val AUTH_HEADER = "Authorization"
    }
}
