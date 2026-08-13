/**
 * Resolve HTTP/WS base URLs for both Vite-dev and the embedded production binary.
 *
 * - Dev: optional VITE_WS_BASE_URL, otherwise ws://localhost:8080/ws (Go API beside Vite).
 * - Production (embedded): same-origin from window.location so LAN IP access works.
 */

function envWsBase(): string {
  const raw = (import.meta.env.VITE_WS_BASE_URL as string | undefined) ?? ''
  return raw.trim()
}

/** WebSocket URL for the player endpoint (.../ws). */
export function getWsBaseUrl(): string {
  const fromEnv = envWsBase()
  if (fromEnv) {
    return fromEnv
  }
  if (import.meta.env.DEV) {
    return 'ws://localhost:8080/ws'
  }
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${proto}//${window.location.host}/ws`
}

/** HTTP origin for REST (admin API), without a trailing slash. */
export function getBaseHttpUrl(): string {
  const fromEnv = envWsBase()
  if (fromEnv) {
    return fromEnv
      .replace(/^wss:/, 'https:')
      .replace(/^ws:/, 'http:')
      .replace(/\/ws$/, '')
  }
  if (import.meta.env.DEV) {
    return 'http://localhost:8080'
  }
  return window.location.origin
}
