import { request } from "@/lib/api/client";

function toQuery(params = {}) {
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") return;
    searchParams.set(key, String(value));
  });

  const query = searchParams.toString();
  return query ? `?${query}` : "";
}

export function listDiscoveryCafes(query = {}) {
  return request(`/discovery/cafes${toQuery(query)}`);
}

export function getDiscoveryCafeById(placeId) {
  return request(`/discovery/cafes/${encodeURIComponent(placeId)}`);
}

export function listCommunityRatingsByPlaceId(placeId) {
  return request(`/community/places/${encodeURIComponent(placeId)}/ratings`);
}
