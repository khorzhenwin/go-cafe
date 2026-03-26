export function formatVisitStatus(status) {
  if (status === "discover") return "Discover";
  if (status === "visited") return "Visited";
  return "Saved";
}

export function formatRating(value) {
  if (!value) return "New";
  return `${Number(value).toFixed(1)} / 5`;
}

export function formatCount(value, singular, plural) {
  return `${value} ${value === 1 ? singular : plural}`;
}

export function formatDate(value) {
  if (!value) return "";

  try {
    return new Intl.DateTimeFormat("en", {
      month: "short",
      day: "numeric",
      year: "numeric"
    }).format(new Date(value));
  } catch {
    return value;
  }
}

export function getCafeSummary(cafe) {
  return cafe.neighborhood || cafe.city || cafe.address || "Coffee spot";
}

export function hasCoordinates(cafe) {
  return typeof cafe?.latitude === "number" && typeof cafe?.longitude === "number";
}
