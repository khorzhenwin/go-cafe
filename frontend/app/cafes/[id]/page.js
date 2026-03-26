"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import AppShell from "@/components/app-shell";
import CafeMap from "@/components/cafe-map";
import ReviewList from "@/components/review-list";
import { useAuth } from "@/components/providers/auth-provider";
import {
  createMyCafe,
  getCafeById,
  getDiscoveryCafeById,
  listCafeRatings,
  listCommunityRatingsByPlaceId,
  listMyCafes
} from "@/lib/api";
import { formatCount, formatRating, formatVisitStatus, getCafeSummary } from "@/lib/presentation";

export default function CafeDetailPage() {
  const params = useParams();
  const cafeId = params?.id;
  const isPersonalCafe = /^\d+$/.test(String(cafeId || ""));
  const { token, ready, isAuthed } = useAuth();
  const [cafe, setCafe] = useState(null);
  const [ratings, setRatings] = useState([]);
  const [myCafes, setMyCafes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  useEffect(() => {
    if (!cafeId) return;

    let cancelled = false;

    async function loadCafe() {
      setLoading(true);
      setError("");

      try {
        const [detail, communityRatings] = await Promise.all(
          isPersonalCafe
            ? [getCafeById(cafeId), listCafeRatings(cafeId)]
            : [getDiscoveryCafeById(cafeId), listCommunityRatingsByPlaceId(cafeId)]
        );

        if (!cancelled) {
          setCafe(detail || null);
          setRatings(communityRatings || []);
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

    loadCafe();

    return () => {
      cancelled = true;
    };
  }, [cafeId, isPersonalCafe]);

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

  const savedCafe = useMemo(() => {
    if (!cafe) return null;

    return (
      myCafes.find((item) => String(item.id) === String(cafe.id)) ||
      myCafes.find((item) => String(item.external_place_id) === String(cafe.external_place_id || cafe.id)) ||
      myCafes.find((item) => String(item.source_cafe_id) === String(cafe.source_cafe_id || cafe.id)) ||
      null
    );
  }, [cafe, myCafes]);

  async function handleSave() {
    if (!cafe) return;

    if (!isAuthed) {
      setMessage("");
      setError("Sign in before saving this cafe.");
      return;
    }

    setSaving(true);
    setError("");
    setMessage("");

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
        visit_status: "to_visit",
        source_provider: cafe.source_provider,
        external_place_id: cafe.external_place_id || cafe.id
      });

      setMyCafes((current) => [created, ...current]);
      setMessage(`${cafe.name} is now in your places.`);
    } catch (saveError) {
      setError(saveError.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <AppShell
      title={cafe?.name || "Cafe detail"}
      subtitle={
        cafe
          ? "Understand the vibe, see community notes, and decide whether this cafe belongs in your personal shortlist."
          : "Loading cafe detail."
      }
      actions={
        <>
          <Link href="/map" className="button button-secondary">
            Back to map
          </Link>
          {savedCafe ? (
            <Link href="/my-places" className="button">
              Open my places
            </Link>
          ) : (
            <button type="button" className="button" onClick={handleSave} disabled={saving || loading}>
              {saving ? "Saving..." : "Save this cafe"}
            </button>
          )}
        </>
      }
    >
      {loading ? <section className="surface empty-state">Loading cafe detail...</section> : null}

      {!loading && error ? <section className="surface empty-state error">{error}</section> : null}

      {!loading && cafe ? (
        <>
          <section className="detail-hero surface">
            <div className="detail-hero-media">
              {cafe.image_url ? (
                <div
                  className="detail-image"
                  aria-hidden="true"
                  style={{ backgroundImage: `url("${cafe.image_url}")` }}
                />
              ) : null}
              <div className="detail-image-fallback" />
            </div>

            <div className="detail-hero-body">
              <div className="badge-row">
                <span className="status-pill">{formatVisitStatus(cafe.visit_status)}</span>
                <span className="ghost-pill">{formatRating(cafe.avg_rating)}</span>
                <span className="ghost-pill">{formatCount(cafe.review_count || 0, "review", "reviews")}</span>
              </div>

              <h2>{cafe.name}</h2>
              <p className="muted">{cafe.address || getCafeSummary(cafe)}</p>
              <p className="body-copy">
                {cafe.description || "This cafe does not have a written note yet, but you can still save it and build your own context."}
              </p>

              <div className="detail-meta">
                <div>
                  <span className="meta-label">City</span>
                  <strong>{cafe.city || "Unknown"}</strong>
                </div>
                <div>
                  <span className="meta-label">Neighborhood</span>
                  <strong>{cafe.neighborhood || "Unknown"}</strong>
                </div>
              </div>

              {message ? <p className="feedback success">{message}</p> : null}
              {error ? <p className="feedback error">{error}</p> : null}
            </div>
          </section>

          <section className="content-grid">
            <CafeMap cafes={[cafe]} selectedCafeId={cafe.id} title="Cafe location" />

            <div className="surface spotlight-card">
              <p className="eyebrow">Next action</p>
              <h2>{savedCafe ? "Already part of your journal" : "Save this place for later"}</h2>
              <p className="body-copy">
                {savedCafe
                  ? "You already have this cafe in your personal collection. Head to My Places to mark it visited or update your note."
                  : "Saving creates a personal copy linked to the original discovery so you can track visit status and write your own review later."}
              </p>
              <div className="card-actions">
                {savedCafe ? (
                  <Link href="/my-places" className="button">
                    Open my places
                  </Link>
                ) : (
                  <button type="button" className="button" onClick={handleSave} disabled={saving}>
                    {saving ? "Saving..." : "Save to my places"}
                  </button>
                )}
                <Link href="/reviews" className="button button-secondary">
                  Review workflow
                </Link>
              </div>
            </div>
          </section>

          <section className="section-stack">
            <div className="section-heading">
              <div>
                <p className="eyebrow">Community reviews</p>
                <h2>How people described this cafe</h2>
              </div>
            </div>
            <ReviewList ratings={ratings} emptyMessage="No community notes yet for this cafe." />
          </section>
        </>
      ) : null}
    </AppShell>
  );
}
