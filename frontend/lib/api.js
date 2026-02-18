const API_BASE_URL = "/api/backend";

export class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {})
    }
  });

  const text = await response.text();
  let payload = null;
  try {
    payload = text ? JSON.parse(text) : null;
  } catch {
    payload = text;
  }

  if (!response.ok) {
    const message = typeof payload === "string" && payload ? payload : `Request failed (${response.status})`;
    throw new ApiError(message, response.status);
  }

  return payload;
}

export function registerUser(body) {
  return request("/auth/register", {
    method: "POST",
    body: JSON.stringify(body)
  });
}

export function loginUser(body) {
  return request("/auth/login", {
    method: "POST",
    body: JSON.stringify(body)
  });
}

export function listMyCafes(token, query = {}) {
  const params = new URLSearchParams();
  if (query.status) params.set("status", query.status);
  if (query.sort) params.set("sort", query.sort);
  const suffix = params.toString() ? `?${params.toString()}` : "";
  return request(`/me/cafes${suffix}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function createMyCafe(token, body) {
  return request("/me/cafes", {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify(body)
  });
}

export function listMyRatings(token) {
  return request("/me/ratings", {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function getCafeById(cafeId) {
  return request(`/cafes/${cafeId}`);
}

export function updateCafe(token, cafeId, body) {
  return request(`/cafes/${cafeId}`, {
    method: "PUT",
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify(body)
  });
}

export function deleteCafe(token, cafeId) {
  return request(`/cafes/${cafeId}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function listCafeRatings(cafeId) {
  return request(`/cafes/${cafeId}/ratings/`);
}

export function createCafeRating(token, cafeId, body) {
  return request(`/cafes/${cafeId}/ratings/`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify(body)
  });
}

export function updateRating(token, ratingId, body) {
  return request(`/ratings/${ratingId}`, {
    method: "PUT",
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify(body)
  });
}

export function deleteRating(token, ratingId) {
  return request(`/ratings/${ratingId}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${token}` }
  });
}
