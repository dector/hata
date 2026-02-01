package hata.core.annotations

import javax.inject.Qualifier


/**
 * Qualifier for compute/CPU-intensive coroutine dispatcher.
 */
@Qualifier
annotation class ComputeDispatcher

/**
 * Qualifier for IO coroutine dispatcher.
 */
@Qualifier
annotation class IoDispatcher
