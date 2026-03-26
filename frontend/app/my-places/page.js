"use client";

import Link from "next/link";
import { useCallback, useEffect, useMemo, useState } from "react";
import AddCafeForm from "@/components/add-cafe-form";
import AppShell from "@/components/app-shell";
import RequireAuth from "@/components/require-auth";
import { useAuth } from "@/components/providers/auth-provider";
import { createMyCafe, deleteCafe, listMyCafes, updateCafe } from "@/lib/api";
import { formatCount, formatVisitStatus } from "@/lib/presentation";

export default function MyPlacesPage() {
  const { token, isAuthed, ready } = useAuth();
  const [cafes, setCafes] = useState([]);
  const [pendingStatusById, setPendingStatusById] = useState({});
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const loadMyCafes = useCallback(async () => {
    if (!ready || !isAuthed) {
      setCafes([]);
      setPendingStatusById({});
      setLoading(false);
      return;
    }

    setLoading(true);
    setError("");

    try {
      const payload = await listMyCafes(token, { sort: "updated_desc" });
      const nextCafes = payload || [];
      setCafes(nextCafes);
      setPendingStatusById(
        nextCafes.reduce((accumulator, cafe) => {
          accumulator[cafe.id] = cafe.visit_status;
          return accumulator;
        }, {})
      );
    } catch (loadError) {
      setError(loadError.message);
    } finally {
      setLoading(false);
    }
  }, [isAuthed, ready, token]);

  useEffect(() => {
    loadMyCafes();
  }, [loadMyCafes]);

  const savedCafes = useMemo(() => cafes.filter((cafe) => cafe.visit_status === "to_visit"), [cafes]);
  const visitedCafes = useMemo(() => cafes.filter((cafe) => cafe.visit_status === "visited"), [cafes]);

  async function handleCreateCafe(payload) {
    setSubmitting(true);
    setError("");
    setMessage("");

    try {
      await createMyCafe(token, payload);
      setMessage("Cafe added to your places.");
      await loadMyCafes();
      return true;
    } catch (createError) {
      setError(createError.message);
      return false;
    } finally {
      setSubmitting(false);
    }
  }

  async function handleStatusUpdate(cafe) {
    setSubmitting(true);
    setError("");
    setMessage("");

    try {
      await updateCafe(token, cafe.id, {
        name: cafe.name,
        address: cafe.address,
        city: cafe.city,
        neighborhood: cafe.neighborhood,
        description: cafe.description,
        image_url: cafe.image_url,
        latitude: cafe.latitude,
        longitude: cafe.longitude,
        visit_status: pendingStatusById[cafe.id] || cafe.visit_status
      });

      setMessage(`${cafe.name} is now marked as ${formatVisitStatus(pendingStatusById[cafe.id])}.`);
      await loadMyCafes();
    } catch (updateError) {
      setError(updateError.message);
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDeleteCafe(cafeId) {
    setSubmitting(true);
    setError("");
    setMessage("");

    try {
      await deleteCafe(token, cafeId);
      setMessage("Cafe removed from your places.");
      await loadMyCafes();
    } catch (deleteError) {
      setError(deleteError.message);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <AppShell
      title="My places"
      subtitle="This is your personal working list: save cafes, mark the ones you visited, and keep your discovery backlog tidy."
      actions={
        <>
          <Link href="/map" className="button button-secondary">
            Discover more cafes
          </Link>
          <Link href="/reviews" className="button button-ghost">
            Write reviews
          </Link>
        </>
      }
    >
      <RequireAuth>
        <section className="stats-grid">
          <article className="surface stat-card">
            <strong>{cafes.length}</strong>
            <span>{formatCount(cafes.length, "place", "places")}</span>
          </article>
          <article className="surface stat-card">
            <strong>{savedCafes.length}</strong>
            <span>Saved for later</span>
          </article>
          <article className="surface stat-card">
            <strong>{visitedCafes.length}</strong>
            <span>Visited</span>
          </article>
        </section>

        <section className="content-grid">
          <AddCafeForm onCreate={handleCreateCafe} submitting={submitting} />

          <section className="surface spotlight-card">
            <p className="eyebrow">Collection notes</p>
            <h2>What belongs here?</h2>
            <p className="body-copy">
              Add a brand-new personal cafe from memory, or save something you discovered elsewhere in the product. Every
              saved place can later be marked visited and reviewed.
            </p>
            {message ? <p className="feedback success">{message}</p> : null}
            {error ? <p className="feedback error">{error}</p> : null}
          </section>
        </section>

        {loading ? <section className="surface empty-state">Loading your places...</section> : null}

        {!loading && !cafes.length ? (
          <section className="surface empty-state">
            Your collection is empty. Add your first cafe or save one from the discovery map.
          </section>
        ) : null}

        {!loading && cafes.length ? (
          <section className="section-stack">
            <div className="section-heading">
              <div>
                <p className="eyebrow">Saved places</p>
                <h2>Places you still want to try</h2>
              </div>
            </div>
            <div className="surface list-surface">
              {savedCafes.length ? (
                savedCafes.map((cafe) => (
                  <article key={cafe.id} className="collection-row">
                    <div>
                      <h3>{cafe.name}</h3>
                      <p className="muted">{cafe.address || cafe.city || "No location summary yet."}</p>
                    </div>
                    <div className="collection-actions">
                      <select
                        value={pendingStatusById[cafe.id] || cafe.visit_status}
                        onChange={(event) =>
                          setPendingStatusById((current) => ({ ...current, [cafe.id]: event.target.value }))
                        }
                      >
                        <option value="to_visit">Saved</option>
                        <option value="visited">Visited</option>
                      </select>
                      <button type="button" className="button button-secondary" onClick={() => handleStatusUpdate(cafe)}>
                        Save status
                      </button>
                      <Link href={`/cafes/${cafe.id}`} className="button button-ghost">
                        View
                      </Link>
                      <button type="button" className="button button-ghost" onClick={() => handleDeleteCafe(cafe.id)}>
                        Remove
                      </button>
                    </div>
                  </article>
                ))
              ) : (
                <p className="muted">No saved cafes yet.</p>
              )}
            </div>

            <div className="section-heading">
              <div>
                <p className="eyebrow">Visited</p>
                <h2>Ready for tasting notes</h2>
              </div>
              <Link href="/reviews" className="button button-secondary">
                Go to review flow
              </Link>
            </div>

            <div className="surface list-surface">
              {visitedCafes.length ? (
                visitedCafes.map((cafe) => (
                  <article key={cafe.id} className="collection-row">
                    <div>
                      <h3>{cafe.name}</h3>
                      <p className="muted">{cafe.address || cafe.city || "No location summary yet."}</p>
                    </div>
                    <div className="collection-actions">
                      <select
                        value={pendingStatusById[cafe.id] || cafe.visit_status}
                        onChange={(event) =>
                          setPendingStatusById((current) => ({ ...current, [cafe.id]: event.target.value }))
                        }
                      >
                        <option value="to_visit">Saved</option>
                        <option value="visited">Visited</option>
                      </select>
                      <button type="button" className="button button-secondary" onClick={() => handleStatusUpdate(cafe)}>
                        Save status
                      </button>
                      <Link href={`/cafes/${cafe.id}`} className="button button-ghost">
                        View
                      </Link>
                      <button type="button" className="button button-ghost" onClick={() => handleDeleteCafe(cafe.id)}>
                        Remove
                      </button>
                    </div>
                  </article>
                ))
              ) : (
                <p className="muted">No visited cafes yet.</p>
              )}
            </div>
          </section>
        ) : null}
      </RequireAuth>
    </AppShell>
  );
}
