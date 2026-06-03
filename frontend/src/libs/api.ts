export const API_BASE_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080";

// Normalize so a trailing slash in VITE_API_URL doesn't yield `//api/...`.
const BASE = API_BASE_URL.replace(/\/+$/, "");

/** Joins a path onto the configured API base URL. */
export function apiUrl(path: string): string {
  return `${BASE}${path}`;
}

/**
 * Fetches a backend endpoint with cookies attached, prefixing the configured
 * API base URL. On a non-2xx response it parses the backend's `{ error }`
 * envelope and throws it as an Error so callers can surface the message.
 */
export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(apiUrl(path), { credentials: "include", ...init });
  const text = await res.text();

  // Tolerate non-JSON bodies (proxy error pages, redirects, empty responses):
  // leave `data` empty on a parse failure so the status-based fallback below
  // produces a clean message instead of a raw SyntaxError.
  let data: any = {};
  try {
    data = text ? JSON.parse(text) : {};
  } catch {
    // ignore — `data` stays {}
  }

  if (!res.ok) {
    throw new Error(data.error ?? `Request failed (${res.status})`);
  }
  return data as T;
}
