"use client";

import Image from "next/image";
import { useMemo } from "react";
import { hasCoordinates } from "@/lib/presentation";

function getCafeKey(cafe) {
  return String(cafe?.external_place_id || cafe?.id || "");
}

function buildStaticMapSrc(cafes, selectedCafe) {
  const params = new URLSearchParams({
    width: "1200",
    height: "720"
  });

  cafes.slice(0, 10).forEach((cafe) => {
    params.append("point", `${cafe.latitude},${cafe.longitude}`);
  });

  if (selectedCafe) {
    params.set("selected", `${selectedCafe.latitude},${selectedCafe.longitude}`);
  }

  return `/api/static-map?${params.toString()}`;
}

export default function CafeMap({ cafes, selectedCafeId, onSelect, title = "Journey map" }) {
  const cafesWithCoordinates = useMemo(() => cafes.filter(hasCoordinates), [cafes]);
  const selectedCafe = useMemo(
    () => cafesWithCoordinates.find((cafe) => getCafeKey(cafe) === String(selectedCafeId)) || cafesWithCoordinates[0] || null,
    [cafesWithCoordinates, selectedCafeId]
  );
  const staticMapSrc = useMemo(
    () => (cafesWithCoordinates.length ? buildStaticMapSrc(cafesWithCoordinates, selectedCafe) : ""),
    [cafesWithCoordinates, selectedCafe]
  );

  return (
    <section className="surface map-panel">
      <div className="cluster">
        <p className="eyebrow">Map view</p>
        <h2>{title}</h2>
        <p className="muted">
          Explore the selected place on a Geoapify static map while using the list below to move between cafes.
        </p>
      </div>

      {cafesWithCoordinates.length ? (
        <div className="map-stage">
          <div className="static-map-image">
            <Image
              src={staticMapSrc}
              alt={selectedCafe ? `${selectedCafe.name} static map` : title}
              fill
              unoptimized
              sizes="(max-width: 1024px) 100vw, 50vw"
            />
          </div>
          <div className="static-map-overlay" aria-hidden="true">
            <span className="static-map-badge">{selectedCafe ? selectedCafe.name : "Discovery view"}</span>
            <span className="static-map-credit">Static map by Geoapify</span>
          </div>
        </div>
      ) : (
        <div className="map-empty">
          <p>No mapped cafes yet.</p>
          <p className="muted">Select a cafe with coordinates to render it on the map.</p>
        </div>
      )}

      {cafesWithCoordinates.length > 1 ? (
        <div className="map-legend">
          {cafesWithCoordinates.slice(0, 8).map((cafe) => {
            const isSelected = getCafeKey(cafe) === String(selectedCafeId);

            return (
              <button
                key={getCafeKey(cafe)}
                type="button"
                className={isSelected ? "map-chip active" : "map-chip"}
                onClick={() => onSelect?.(cafe)}
              >
                {cafe.name}
              </button>
            );
          })}
        </div>
      ) : null}
    </section>
  );
}
