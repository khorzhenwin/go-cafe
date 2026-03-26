"use client";

import { buildGoogleMapsEmbedUrl } from "@/lib/presentation";

export default function CafeMap({ cafes, selectedCafeId, onSelect, title = "Journey map" }) {
  const selectedCafe =
    cafes.find((cafe) => String(cafe.id) === String(selectedCafeId)) ||
    cafes.find((cafe) => String(cafe.external_place_id) === String(selectedCafeId)) ||
    cafes[0] ||
    null;
  const embedUrl = selectedCafe ? buildGoogleMapsEmbedUrl(selectedCafe) : "";

  return (
    <section className="surface map-panel">
      <div className="cluster">
        <p className="eyebrow">Map view</p>
        <h2>{title}</h2>
        <p className="muted">
          Explore the selected place on an actual Google map embed instead of the old placeholder canvas.
        </p>
      </div>

      {selectedCafe && embedUrl ? (
        <div className="map-stage">
          <iframe
            key={embedUrl}
            title={`${selectedCafe.name} map`}
            src={embedUrl}
            loading="lazy"
            referrerPolicy="no-referrer-when-downgrade"
            className="map-embed"
          />
        </div>
      ) : (
        <div className="map-empty">
          <p>No mapped cafes yet.</p>
          <p className="muted">Select a cafe with an address or coordinates to render it on Google Maps.</p>
        </div>
      )}

      {cafes.length > 1 ? (
        <div className="map-legend">
          {cafes.slice(0, 8).map((cafe) => {
            const isSelected =
              String(cafe.id) === String(selectedCafeId) || String(cafe.external_place_id) === String(selectedCafeId);

            return (
              <button
                key={cafe.id || cafe.external_place_id}
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
