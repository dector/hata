package hata.core.annotations

import javax.inject.Qualifier

/**
 * Qualifier for compute/CPU-intensive coroutine dispatcher.
 *
 * This annotation marks CoroutineDispatcher instances intended for CPU-intensive operations.
 * Use this qualifier to inject the compute dispatcher in classes that need it.
 *
 * Example usage in Hilt modules:
 * ```
 * @Provides
 * @ComputeDispatcher
 * fun computeDispatcher(): CoroutineDispatcher = Dispatchers.Default
 * ```
 *
 * Example usage in non-Hilt contexts:
 * ```
 * // Use as documentation/marker without DI
 * class MyService(
 *     @ComputeDispatcher private val dispatcher: CoroutineDispatcher = Dispatchers.Default
 * )
 * ```
 */
@Qualifier
annotation class ComputeDispatcher

/**
 * Qualifier for IO coroutine dispatcher.
 *
 * This annotation marks CoroutineDispatcher instances intended for IO operations.
 * Use this qualifier to inject the IO dispatcher in classes that need it.
 *
 * Example usage in Hilt modules:
 * ```
 * @Provides
 * @IoDispatcher
 * fun ioDispatcher(): CoroutineDispatcher = Dispatchers.IO
 * ```
 *
 * Example usage in non-Hilt contexts:
 * ```
 * // Use as documentation/marker without DI
 * class MyService(
 *     @IoDispatcher private val dispatcher: CoroutineDispatcher = Dispatchers.IO
 * )
 * ```
 */
@Qualifier
annotation class IoDispatcher
