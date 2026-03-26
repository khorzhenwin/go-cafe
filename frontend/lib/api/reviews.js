import { authHeaders, request } from "@/lib/api/client";

export function listMyRatings(token) {
  return request("/me/ratings", {
    headers: authHeaders(token)
  });
}

export function listCafeRatings(cafeId) {
  return request(`/cafes/${cafeId}/ratings/`);
}

export function createCafeRating(token, cafeId, body) {
  return request(`/cafes/${cafeId}/ratings/`, {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify(body)
  });
}

export function updateRating(token, ratingId, body) {
  return request(`/ratings/${ratingId}`, {
    method: "PUT",
    headers: authHeaders(token),
    body: JSON.stringify(body)
  });
}

export function deleteRating(token, ratingId) {
  return request(`/ratings/${ratingId}`, {
    method: "DELETE",
    headers: authHeaders(token)
  });
}
