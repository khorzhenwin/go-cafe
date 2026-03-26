"use client";

import { formatDate } from "@/lib/presentation";

export default function ReviewList({ ratings, emptyMessage = "No reviews yet.", onDelete, canDelete = false }) {
  if (!ratings.length) {
    return <section className="surface empty-state">{emptyMessage}</section>;
  }

  return (
    <section className="review-list">
      {ratings.map((rating) => (
        <article key={rating.id} className="surface review-card">
          <div className="cluster-sm">
            <div>
              <p className="eyebrow">{rating.cafe_listing?.name || `Cafe #${rating.cafe_listing_id}`}</p>
              <h3>{rating.rating}/5</h3>
            </div>
            <div className="review-card-meta">
              <span>{rating.user?.name || "Cafe Explorer"}</span>
              <span>{formatDate(rating.visited_at)}</span>
            </div>
          </div>

          <p className="body-copy">{rating.review || "No written tasting note yet."}</p>

          {canDelete && onDelete ? (
            <button type="button" className="button button-ghost" onClick={() => onDelete(rating.id)}>
              Delete review
            </button>
          ) : null}
        </article>
      ))}
    </section>
  );
}
