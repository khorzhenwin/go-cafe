const API_BASE_URL = "/api/backend";

export class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export async function request(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {})
    },
    cache: "no-store"
  });

  const text = await response.text();
  let payload = null;

  try {
    payload = text ? JSON.parse(text) : null;
  } catch {
    payload = text;
  }

  if (!response.ok) {
    if (response.status === 401 && typeof window !== "undefined") {
      window.localStorage.removeItem("go-cafe-token");
      window.dispatchEvent(new Event("go-cafe-auth-changed"));
    }

    const message =
      typeof payload === "string"
        ? payload
        : payload?.detail || payload?.error || payload?.message || `Request failed (${response.status})`;
    throw new ApiError(message, response.status);
  }

  return payload;
}

export function authHeaders(token) {
  return token ? { Authorization: `Bearer ${token}` } : {};
}
