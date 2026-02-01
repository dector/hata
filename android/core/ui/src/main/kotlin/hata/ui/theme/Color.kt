package hata.ui.theme

import androidx.compose.ui.graphics.Color

/**
 * Hata app color palette following Material Design 3 principles.
 *
 * Core colors are organized by semantic role:
 * - Primary: Brand identity and key actions
 * - Surface: Backgrounds and containers
 * - Content: Text and icons
 * - State: Interactive and feedback colors
 *
 * Usage: Access colors via `HataColors.primary`, `HataColors.surface`, etc.
 * Theme integration: Colors are automatically available via MaterialTheme.colorScheme
 */
object HataColors {
    // Primary brand colors
    val primary = Color(0xFF8FB899)              // Primary green accent
    val onPrimary = Color(0xFF1E2423)             // Text/icons on primary

    // Surface colors (backgrounds & cards)
    val background = Color(0xFF1E2423)            // Main app background
    val surface = Color(0xFF2D3533)               // Card backgrounds
    val surfaceVariant = Color(0xFF3D4845)        // Icon container backgrounds

    // Content colors (text & icons)
    val onSurface = Color.White                   // Primary text
    val onSurfaceVariant = Color(0xFF6B7875)      // Secondary text/inactive states

    // State colors (interactive elements)
    val primaryContainer = Color(0xFFA8C5B0)      // Active state backgrounds
    val primaryContainerVariant = Color(0xFFC8DFD0) // Switch track (checked)

    // Status colors (feedback & alerts)
    val error = Color(0xFFE57373)                 // Error messages
    val errorVariant = Color(0xFFEF5350)          // Notification indicator
}

/**
 * Accent colors for specific features that don't fit Material Design roles.
 */
object HataAccentColors {
    val avatarBackground = Color(0xFFB39B8D)      // Avatar background tan
    val avatarIcon = Color(0xFFE8D4C4)            // Avatar icon tint
}

// Material Design template colors - kept for theme compatibility
internal val Purple80 = Color(0xFFD0BCFF)
internal val PurpleGrey80 = Color(0xFFCCC2DC)
internal val Pink80 = Color(0xFFEFB8C8)

internal val Purple40 = Color(0xFF6650a4)
internal val PurpleGrey40 = Color(0xFF625b71)
internal val Pink40 = Color(0xFF7D5260)
