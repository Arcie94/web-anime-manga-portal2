// Use backend service name for SSR, public URL for client
const isServer = typeof window === 'undefined';
// @ts-ignore
export const API_BASE = isServer ? (process.env.API_URL || "http://localhost:3000/api") : "/api";
console.log("DEBUG: API_BASE =", API_BASE, "isServer =", isServer);

type FetchOptions = {
  method?: string;
  body?: any;
  headers?: Record<string, string>;
};

export async function apiFetch<T>(
  endpoint: string,
  options: FetchOptions = {}
): Promise<T> {
  const url = `${API_BASE}${endpoint}`;
  if (isServer) console.log(`Fetching: ${url}`);
  const headers = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  const config: RequestInit = {
    method: options.method || "GET",
    headers,
  };

  if (options.body) {
    config.body = JSON.stringify(options.body);
  }

  const res = await fetch(url, config);
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(
      errorData.error || `Request failed with status ${res.status}`
    );
  }

  return res.json();
}
