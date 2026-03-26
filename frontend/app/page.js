"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import AppShell from "@/components/app-shell";
import CafeCard from "@/components/cafe-card";
import CafeMap from "@/components/cafe-map";
import { useAuth } from "@/components/providers/auth-provider";
import { listDiscoveryCafes } from "@/lib/api";

const BENEFITS = [
  "Browse real cafes from Google instead of relying on seeded shared database content.",
  "Save interesting places into a personal shortlist for later visits.",
  "Turn visits into reviews without losing the discovery context."
];

export default function HomePage() {
  const { isAuthed } = useAuth();
  const [cafes, setCafes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedCafe, setSelectedCafe] = useState(null);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      setLoading(true);
      setError("");

      try {
        const payload = await listDiscoveryCafes({ city: "Singapore", limit: 6 });
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

    load();

    return () => {
      cancelled = true;
    };
  }, []);

  const stats = useMemo(
    () => [
      { label: "Community spots", value: cafes.length },
      { label: "Mapped entries", value: cafes.filter((cafe) => cafe.latitude && cafe.longitude).length },
      { label: "Total reviews", value: cafes.reduce((sum, cafe) => sum + (cafe.review_count || 0), 0) }
    ],
    [cafes]
  );

  return (
    <AppShell
      title="Discover cafes with better context"
      subtitle="Cafe Hub now leads with discovery: browse Google-sourced cafes, understand the vibe quickly, then save the ones worth adding to your own ritual."
      actions={
        <>
          <Link href="/map" className="button">
            Explore the map
          </Link>
          <Link href={isAuthed ? "/my-places" : "/auth"} className="button button-secondary">
            {isAuthed ? "Open my places" : "Create an account"}
          </Link>
        </>
      }
    >
      <section className="hero-grid">
        <div className="surface hero-card">
          <p className="eyebrow">Why this redesign matters</p>
          <h2>Intent is clearer from the first screen</h2>
          <div className="benefit-list">
            {BENEFITS.map((benefit) => (
              <p key={benefit}>{benefit}</p>
            ))}
          </div>
          <div className="hero-actions">
            <Link href="/map" className="button">
              Start discovering
            </Link>
            <Link href="/reviews" className="button button-ghost">
              See the journal flow
            </Link>
          </div>
        </div>

        <div className="stats-grid">
          {stats.map((item) => (
            <article key={item.label} className="surface stat-card">
              <strong>{item.value}</strong>
              <span>{item.label}</span>
            </article>
          ))}
          <article className="surface stat-card stat-card-highlight">
            <strong>{cafes.length}</strong>
            <span>Featured cafes right now</span>
          </article>
        </div>
      </section>

      <section className="content-grid">
        <CafeMap
          cafes={cafes}
          selectedCafeId={selectedCafe?.external_place_id || selectedCafe?.id}
          onSelect={setSelectedCafe}
          title="Featured discovery map"
        />

        <div className="surface spotlight-card">
          <p className="eyebrow">Selected cafe</p>
          {selectedCafe ? (
            <>
              <h2>{selectedCafe.name}</h2>
              <p className="muted">{selectedCafe.address || selectedCafe.city || "Google Places result"}</p>
              <p className="body-copy">
                {selectedCafe.description || "Open this cafe detail to see the full context and community notes."}
              </p>
              <div className="card-actions">
                <Link href={`/cafes/${encodeURIComponent(selectedCafe.id)}`} className="button">
                  Open cafe detail
                </Link>
                <Link href="/map" className="button button-secondary">
                  Browse all cafes
                </Link>
              </div>
            </>
          ) : loading ? (
            <p>Loading featured cafes...</p>
          ) : (
            <p>No public cafes yet. Add one from your personal area to populate the map.</p>
          )}
          {error ? <p className="feedback error">{error}</p> : null}
        </div>
      </section>

      <section className="section-stack">
        <div className="section-heading">
          <div>
            <p className="eyebrow">Featured cafes</p>
            <h2>High-signal places to start with</h2>
          </div>
          <Link href="/map" className="button button-secondary">
            View full discovery list
          </Link>
        </div>

        {loading ? <section className="surface empty-state">Loading community picks...</section> : null}

        {!loading && !cafes.length ? (
          <section className="surface empty-state">
            Google discovery is not available right now. Set the Google Places API key to populate this view with live cafes.
          </section>
        ) : null}

        {!loading ? (
          <div className="card-grid">
            {cafes.map((cafe) => (
              <CafeCard key={cafe.id} cafe={cafe} />
            ))}
          </div>
        ) : null}
      </section>
    </AppShell>
  );
}
