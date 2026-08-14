/**
 * Resolve HTTP/WS base URLs for both Vite-dev and the embedded production binary.
 *
 * - Dev: optional VITE_WS_BASE_URL, otherwise ws://localhost:8080/ws (Go API beside Vite).
 * - Production (embedded): always same-origin from window.location so LAN IP access works.
 *   VITE_WS_BASE_URL is intentionally ignored in production — baking localhost into the
 *   binary would break every student browser on the classroom LAN.
 */

function envWsBase(): string {
  const raw = (import.meta.env.VITE_WS_BASE_URL as string | undefined) ?? ''
  return raw.trim()
}

/** WebSocket URL for the player endpoint (.../ws). */
export function getWsBaseUrl(): string {
  // Production builds must never honor a bake-time override.
  if (!import.meta.env.DEV) {
    const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    return `${proto}//${window.location.host}/ws`
  }
  const fromEnv = envWsBase()
  if (fromEnv) {
    return fromEnv
  }
  return 'ws://localhost:8080/ws'
}

/** HTTP origin for REST (admin API), without a trailing slash. */
export function getBaseHttpUrl(): string {
  if (!import.meta.env.DEV) {
    return window.location.origin
  }
  const fromEnv = envWsBase()
  if (fromEnv) {
    return fromEnv
      .replace(/^wss:/, 'https:')
      .replace(/^ws:/, 'http:')
      .replace(/\/ws$/, '')
  }
  return 'http://localhost:8080'
}
