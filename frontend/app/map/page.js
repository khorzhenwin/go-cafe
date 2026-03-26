"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import AppShell from "@/components/app-shell";
import CafeCard from "@/components/cafe-card";
import CafeMap from "@/components/cafe-map";
import { useAuth } from "@/components/providers/auth-provider";
import { createMyCafe, listDiscoveryCafes, listMyCafes } from "@/lib/api";

export default function MapPage() {
  const { token, isAuthed, ready } = useAuth();
  const [query, setQuery] = useState("");
  const [city, setCity] = useState("Singapore");
  const [cafes, setCafes] = useState([]);
  const [myCafes, setMyCafes] = useState([]);
  const [selectedCafe, setSelectedCafe] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadDiscovery() {
      setLoading(true);
      setError("");

      try {
        const payload = await listDiscoveryCafes({ query, city, limit: 24 });
        if (!cancelled) {
          setCafes(payload || []);
          setSelectedCafe(payload?.[0] || null);
        }
      } catch (loadError) {
        if (!cancelled) {
          setError(loadError.message);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadDiscovery();

    return () => {
      cancelled = true;
    };
  }, [city, query]);

  useEffect(() => {
    if (!ready || !isAuthed) {
      setMyCafes([]);
      return;
    }

    let cancelled = false;

    async function loadMyCafes() {
      try {
        const payload = await listMyCafes(token, { sort: "updated_desc" });
        if (!cancelled) {
          setMyCafes(payload || []);
        }
      } catch {
        if (!cancelled) {
          setMyCafes([]);
        }
      }
    }

    loadMyCafes();

    return () => {
      cancelled = true;
    };
  }, [isAuthed, ready, token]);

  const savedMap = useMemo(() => {
    return new Map(
      myCafes.map((cafe) => {
        const key = cafe.external_place_id || cafe.source_cafe_id || cafe.id;
        return [String(key), cafe];
      })
    );
  }, [myCafes]);

  async function handleSaveCafe(cafe) {
    if (!isAuthed) {
      setMessage("");
      setError("Sign in before saving a cafe to your personal collection.");
      return;
    }

    setSaving(true);
    setMessage("");
    setError("");

    try {
      const created = await createMyCafe(token, {
        name: cafe.name,
        address: cafe.address,
        city: cafe.city,
        neighborhood: cafe.neighborhood,
        description: cafe.description,
        image_url: cafe.image_url,
        latitude: cafe.latitude,
        longitude: cafe.longitude,
        source_provider: cafe.source_provider,
        external_place_id: cafe.external_place_id || cafe.id,
        visit_status: "to_visit"
      });

      setMyCafes((current) => [created, ...current]);
      setMessage(`${cafe.name} was added to your places.`);
    } catch (saveError) {
      setError(saveError.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <AppShell
      title="Map-first cafe discovery"
      subtitle="Search real cafe results from Google by name, neighborhood, or city. Save the ones worth revisiting into your own private journal."
      actions={
        <>
          <Link href="/my-places" className="button button-secondary">
            Open my places
          </Link>
          <Link href="/reviews" className="button button-ghost">
            Review flow
          </Link>
        </>
      }
    >
      <section className="surface filter-bar">
        <label>
          Search
          <input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search by cafe name, address, or vibe"
          />
        </label>

        <label>
          City
          <input value={city} onChange={(event) => setCity(event.target.value)} placeholder="Singapore" />
        </label>
      </section>

      <section className="content-grid">
        <CafeMap
          cafes={cafes}
          selectedCafeId={selectedCafe?.external_place_id || selectedCafe?.id}
          onSelect={setSelectedCafe}
          title="Discovery map"
        />

        <div className="surface spotlight-card">
          <p className="eyebrow">Selected place</p>
          {selectedCafe ? (
            <>
              <h2>{selectedCafe.name}</h2>
              <p className="muted">{selectedCafe.address || "Google discovery result"}</p>
              <p className="body-copy">
                {selectedCafe.description || "Open the detail page to understand the context behind this recommendation."}
              </p>
              <div className="card-actions">
                <Link href={`/cafes/${encodeURIComponent(selectedCafe.id)}`} className="button">
                  Open detail
                </Link>
                {savedMap.has(String(selectedCafe.external_place_id || selectedCafe.id)) ? (
                  <Link href="/my-places" className="button button-secondary">
                    Already in my places
                  </Link>
                ) : (
                  <button type="button" className="button button-secondary" onClick={() => handleSaveCafe(selectedCafe)}>
                    {saving ? "Saving..." : "Save to my places"}
                  </button>
                )}
              </div>
            </>
          ) : loading ? (
            <p>Loading discovery map...</p>
          ) : (
            <p>No cafes match your current filters.</p>
          )}
          {message ? <p className="feedback success">{message}</p> : null}
          {error ? <p className="feedback error">{error}</p> : null}
        </div>
      </section>

      {loading ? <section className="surface empty-state">Loading discovery results...</section> : null}

      {!loading && !cafes.length ? (
        <section className="surface empty-state">
          No Google cafe results match your current filters. Try broadening the search or check that the Google Places API key is configured.
        </section>
      ) : null}

      {!loading ? (
        <div className="card-grid">
          {cafes.map((cafe) => {
            const savedCafe = savedMap.get(String(cafe.external_place_id || cafe.id));

            return (
              <CafeCard
                key={cafe.id}
                cafe={cafe}
                action={
                  savedCafe ? (
                    <Link href="/my-places" className="button button-ghost">
                      Saved already
                    </Link>
                  ) : (
                    <button type="button" className="button button-ghost" onClick={() => handleSaveCafe(cafe)}>
                      Save
                    </button>
                  )
                }
              />
            );
          })}
        </div>
      ) : null}
    </AppShell>
  );
}
