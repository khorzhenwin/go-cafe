"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import AppShell from "@/components/app-shell";
import RequireAuth from "@/components/require-auth";
import ReviewList from "@/components/review-list";
import { useAuth } from "@/components/providers/auth-provider";
import { createCafeRating, deleteRating, listMyCafes, listMyRatings } from "@/lib/api";

export default function ReviewsPage() {
  const { token, isAuthed, ready } = useAuth();
  const [cafes, setCafes] = useState([]);
  const [ratings, setRatings] = useState([]);
  const [selectedCafeId, setSelectedCafeId] = useState("");
  const [score, setScore] = useState("5");
  const [review, setReview] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const loadData = useCallback(async () => {
    if (!ready || !isAuthed) {
      setCafes([]);
      setRatings([]);
      setSelectedCafeId("");
      setLoading(false);
      return;
    }

    setLoading(true);
    setError("");

    try {
      const [nextCafes, nextRatings] = await Promise.all([
        listMyCafes(token, { sort: "updated_desc" }),
        listMyRatings(token)
      ]);

      const cafesPayload = nextCafes || [];
      const visitedOnly = cafesPayload.filter((cafe) => cafe.visit_status === "visited");

      setCafes(cafesPayload);
      setRatings(nextRatings || []);
      setSelectedCafeId((current) => current || String(visitedOnly[0]?.id || ""));
    } catch (loadError) {
      setError(loadError.message);
    } finally {
      setLoading(false);
    }
  }, [isAuthed, ready, token]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const visitedCafes = useMemo(() => cafes.filter((cafe) => cafe.visit_status === "visited"), [cafes]);

  async function handleCreateReview(event) {
    event.preventDefault();
    setSubmitting(true);
    setError("");
    setMessage("");

    try {
      await createCafeRating(token, selectedCafeId, {
        visited_at: new Date().toISOString(),
        rating: Number(score),
        review
      });

      setReview("");
      setScore("5");
      setMessage("Review saved.");
      await loadData();
    } catch (submitError) {
      setError(submitError.message);
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDeleteReview(ratingId) {
    setSubmitting(true);
    setError("");
    setMessage("");

    try {
      await deleteRating(token, ratingId);
      setMessage("Review deleted.");
      await loadData();
    } catch (deleteError) {
      setError(deleteError.message);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <AppShell
      title="Your tasting notes"
      subtitle="Reviews now live in their own space so you can focus on writing, revisiting, and pruning journal entries without juggling the discovery interface."
    >
      <RequireAuth>
        <section className="content-grid">
          <section className="surface">
            <div className="cluster">
              <p className="eyebrow">Write a review</p>
              <h2>Capture the visit while it is fresh</h2>
            </div>

            <form className="stack-form" onSubmit={handleCreateReview}>
              <label>
                Visited cafe
                <select value={selectedCafeId} onChange={(event) => setSelectedCafeId(event.target.value)}>
                  <option value="">Select a visited cafe</option>
                  {visitedCafes.map((cafe) => (
                    <option key={cafe.id} value={cafe.id}>
                      {cafe.name}
                    </option>
                  ))}
                </select>
              </label>

              <label>
                Score
                <select value={score} onChange={(event) => setScore(event.target.value)}>
                  {[5, 4, 3, 2, 1].map((value) => (
                    <option key={value} value={value}>
                      {value}/5
                    </option>
                  ))}
                </select>
              </label>

              <label>
                Tasting note
                <textarea
                  value={review}
                  onChange={(event) => setReview(event.target.value)}
                  placeholder="What stood out? Espresso texture, lighting, pace, or the feeling of the room?"
                />
              </label>

              <button type="submit" className="button" disabled={submitting || !selectedCafeId || !visitedCafes.length}>
                {submitting ? "Saving..." : "Save review"}
              </button>
            </form>

            {message ? <p className="feedback success">{message}</p> : null}
            {error ? <p className="feedback error">{error}</p> : null}
            {!visitedCafes.length && !loading ? (
              <p className="muted">Mark a cafe as visited from My Places before writing a review.</p>
            ) : null}
          </section>

          <section className="surface spotlight-card">
            <p className="eyebrow">Review cadence</p>
            <h2>Separate writing from discovery</h2>
            <p className="body-copy">
              The new flow keeps exploration on the map, collection management in My Places, and your reflective writing here.
            </p>
          </section>
        </section>

        {loading ? <section className="surface empty-state">Loading your reviews...</section> : null}

        {!loading ? (
          <section className="section-stack">
            <div className="section-heading">
              <div>
                <p className="eyebrow">Review history</p>
                <h2>Recent journal entries</h2>
              </div>
            </div>
            <ReviewList
              ratings={ratings}
              emptyMessage="You have not written any reviews yet."
              onDelete={handleDeleteReview}
              canDelete
            />
          </section>
        ) : null}
      </RequireAuth>
    </AppShell>
  );
}
