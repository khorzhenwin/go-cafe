"use client";

import Link from "next/link";
import { formatCount, formatRating, formatVisitStatus, getCafeSummary } from "@/lib/presentation";

export default function CafeCard({
  cafe,
  href = `/cafes/${encodeURIComponent(cafe.external_place_id || cafe.id)}`,
  action,
  secondaryAction,
  variant = "discovery"
}) {
  return (
    <article className={`surface cafe-card cafe-card-${variant}`}>
      <div className="cafe-card-media">
        {cafe.image_url ? (
          <div
            className="cafe-card-image"
            aria-hidden="true"
            style={{ backgroundImage: `url("${cafe.image_url}")` }}
          />
        ) : null}
        <div className="cafe-card-gradient" />
        <div className="cafe-card-badges">
          <span className="status-pill">{formatVisitStatus(cafe.visit_status)}</span>
          {cafe.city ? <span className="ghost-pill">{cafe.city}</span> : null}
        </div>
      </div>

      <div className="cafe-card-body">
        <div className="cluster-sm">
          <div>
            <h3>{cafe.name}</h3>
            <p className="muted">{getCafeSummary(cafe)}</p>
          </div>
          <div className="rating-meta">
            <strong>{formatRating(cafe.avg_rating)}</strong>
            <span>{formatCount(cafe.review_count || 0, "review", "reviews")}</span>
          </div>
        </div>

        <p className="body-copy">{cafe.description || "A real-world cafe result worth keeping on your radar."}</p>

        <div className="card-actions">
          <Link href={href} className="button button-secondary">
            View details
          </Link>
          {action}
          {secondaryAction}
        </div>
      </div>
    </article>
  );
}
