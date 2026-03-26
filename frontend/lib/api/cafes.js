import { authHeaders, request } from "@/lib/api/client";

function toQuery(params = {}) {
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") return;
    searchParams.set(key, String(value));
  });

  const query = searchParams.toString();
  return query ? `?${query}` : "";
}

export function getCafeById(cafeId) {
  return request(`/cafes/${cafeId}`);
}

export function listMyCafes(token, query = {}) {
  return request(`/me/cafes${toQuery(query)}`, {
    headers: authHeaders(token)
  });
}

export function createMyCafe(token, body) {
  return request("/me/cafes", {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify(body)
  });
}

export function updateCafe(token, cafeId, body) {
  return request(`/cafes/${cafeId}`, {
    method: "PUT",
    headers: authHeaders(token),
    body: JSON.stringify(body)
  });
}

export function deleteCafe(token, cafeId) {
  return request(`/cafes/${cafeId}`, {
    method: "DELETE",
    headers: authHeaders(token)
  });
}

export function searchCafeAddresses(query, limit = 5) {
  return request(`/cafes/autocomplete${toQuery({ text: query, limit })}`);
}
