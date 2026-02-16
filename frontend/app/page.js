"use client";

import { useEffect, useMemo, useState } from "react";
import {
  createCafeRating,
  createMyCafe,
  deleteCafe,
  deleteRating,
  listCafeRatings,
  listMyCafes,
  listMyRatings,
  loginUser,
  registerUser,
  updateCafe,
  updateRating
} from "@/lib/api";

export default function HomePage() {
  const [authMode, setAuthMode] = useState("signup");
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [token, setToken] = useState("");
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const [cafes, setCafes] = useState([]);
  const [ratings, setRatings] = useState([]);
  const [ratingByCafe, setRatingByCafe] = useState({});

  const [cafeName, setCafeName] = useState("");
  const [cafeAddress, setCafeAddress] = useState("");
  const [cafeDescription, setCafeDescription] = useState("");

  const [selectedCafeId, setSelectedCafeId] = useState("");
  const [reviewRating, setReviewRating] = useState("5");
  const [reviewText, setReviewText] = useState("");

  const [editCafeId, setEditCafeId] = useState("");
  const [editCafeName, setEditCafeName] = useState("");
  const [editCafeAddress, setEditCafeAddress] = useState("");
  const [editCafeDescription, setEditCafeDescription] = useState("");

  const [editRatingId, setEditRatingId] = useState("");
  const [editRatingScore, setEditRatingScore] = useState("5");
  const [editRatingReview, setEditRatingReview] = useState("");
  const [sortBy, setSortBy] = useState("updated_desc");

  const isAuthed = useMemo(() => token.trim().length > 0, [token]);

  useEffect(() => {
    const saved = window.localStorage.getItem("go-cafe-token");
    if (saved) setToken(saved);
  }, []);

  useEffect(() => {
    if (token) window.localStorage.setItem("go-cafe-token", token);
    else window.localStorage.removeItem("go-cafe-token");
  }, [token]);

  const cafesWithMeta = useMemo(() => {
    const withMeta = cafes.map((cafe) => {
      const cafeRatings = ratingByCafe[cafe.id] || [];
      const count = cafeRatings.length;
      const avg = count ? cafeRatings.reduce((acc, r) => acc + Number(r.rating || 0), 0) / count : 0;
      return { ...cafe, rating_count: count, rating_avg: avg };
    });
    return withMeta.sort((a, b) => {
      if (sortBy === "name_asc") return a.name.localeCompare(b.name);
      if (sortBy === "name_desc") return b.name.localeCompare(a.name);
      if (sortBy === "rating_desc") return b.rating_avg - a.rating_avg;
      if (sortBy === "created_desc") return new Date(b.created_at) - new Date(a.created_at);
      return new Date(b.updated_at) - new Date(a.updated_at);
    });
  }, [cafes, ratingByCafe, sortBy]);

  async function handleRegister() {
    setError("");
    setMessage("");
    setLoading(true);
    try {
      const response = await registerUser({ email, name, password });
      setToken(response.token || "");
      setMessage("Signed up successfully.");
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleLogin() {
    setError("");
    setMessage("");
    setLoading(true);
    try {
      const response = await loginUser({ email, password });
      setToken(response.token || "");
      setMessage("Logged in successfully.");
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleRefreshData() {
    if (!isAuthed) return;
    setError("");
    setMessage("");
    setLoading(true);
    try {
      const [myCafes, myRatings] = await Promise.all([listMyCafes(token), listMyRatings(token)]);
      const cafesData = myCafes || [];
      setCafes(cafesData);
      setRatings(myRatings || []);
      const ratingsByCafeEntries = await Promise.all(
        cafesData.map(async (cafe) => {
          const cafeRatings = await listCafeRatings(cafe.id);
          return [cafe.id, cafeRatings || []];
        })
      );
      setRatingByCafe(Object.fromEntries(ratingsByCafeEntries));
      setSelectedCafeId((prev) => prev || String(cafesData[0]?.id || ""));
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleCreateCafe() {
    if (!isAuthed) return;
    setError("");
    setMessage("");
    setLoading(true);
    try {
      await createMyCafe(token, { name: cafeName, address: cafeAddress, description: cafeDescription });
      setCafeName("");
      setCafeAddress("");
      setCafeDescription("");
      setMessage("Cafe listing created.");
      await handleRefreshData();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleUpdateCafe() {
    if (!editCafeId) return;
    setError("");
    setMessage("");
    setLoading(true);
    try {
      await updateCafe(token, editCafeId, {
        name: editCafeName,
        address: editCafeAddress,
        description: editCafeDescription
      });
      setMessage("Cafe listing updated.");
      await handleRefreshData();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleDeleteCafe(cafeId) {
    setError("");
    setMessage("");
    setLoading(true);
    try {
      await deleteCafe(token, cafeId);
      if (String(cafeId) === String(selectedCafeId)) setSelectedCafeId("");
      setMessage("Cafe listing deleted.");
      await handleRefreshData();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleCreateReview() {
    if (!selectedCafeId) return;
    setError("");
    setMessage("");
    setLoading(true);
    try {
      await createCafeRating(token, selectedCafeId, {
        visited_at: new Date().toISOString(),
        rating: Number(reviewRating),
        review: reviewText
      });
      setReviewText("");
      setReviewRating("5");
      setMessage("Review added.");
      await handleRefreshData();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  function startCafeEdit(cafe) {
    setEditCafeId(String(cafe.id));
    setEditCafeName(cafe.name || "");
    setEditCafeAddress(cafe.address || "");
    setEditCafeDescription(cafe.description || "");
  }

  function startRatingEdit(rating) {
    setEditRatingId(String(rating.id));
    setEditRatingScore(String(rating.rating || 5));
    setEditRatingReview(rating.review || "");
  }

  async function handleUpdateRating() {
    if (!editRatingId) return;
    const original = ratings.find((r) => String(r.id) === String(editRatingId));
    setError("");
    setMessage("");
    setLoading(true);
    try {
      await updateRating(token, editRatingId, {
        visited_at: original?.visited_at || new Date().toISOString(),
        rating: Number(editRatingScore),
        review: editRatingReview
      });
      setMessage("Review updated.");
      await handleRefreshData();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleDeleteRating(ratingId) {
    setError("");
    setMessage("");
    setLoading(true);
    try {
      await deleteRating(token, ratingId);
      setMessage("Review deleted.");
      await handleRefreshData();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  function handleLogout() {
    setToken("");
    setCafes([]);
    setRatings([]);
    setRatingByCafe({});
    setMessage("Logged out.");
    setError("");
  }

  async function handleAuthSubmit(e) {
    e.preventDefault();
    if (authMode === "signup") await handleRegister();
    else await handleLogin();
  }

  return (
    <main className="app-shell">
      <header className="hero">
        <h1>Cafe Hub</h1>
        <p className="muted">Sign up, create cafe listings, write reviews, and sort your collection.</p>
      </header>

      {!isAuthed ? (
        <section className="card auth-card">
          <div className="auth-tabs">
            <button className={authMode === "signup" ? "tab active" : "tab"} onClick={() => setAuthMode("signup")}>
              Sign up
            </button>
            <button className={authMode === "login" ? "tab active" : "tab"} onClick={() => setAuthMode("login")}>
              Login
            </button>
          </div>
          <form onSubmit={handleAuthSubmit}>
            <div className="row">
              <div>
                <label htmlFor="email">Email</label>
                <input id="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
              </div>
              {authMode === "signup" ? (
                <div>
                  <label htmlFor="name">Display name</label>
                  <input id="name" value={name} onChange={(e) => setName(e.target.value)} required />
                </div>
              ) : null}
              <div>
                <label htmlFor="password">Password</label>
                <input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
            </div>
            <div className="actions">
              <button type="submit" disabled={loading}>
                {loading ? "Please wait..." : authMode === "signup" ? "Create account" : "Login"}
              </button>
            </div>
          </form>
        </section>
      ) : (
        <>
          <section className="card">
            <div className="toolbar">
              <h2>Workspace</h2>
              <div className="actions">
                <button className="secondary" onClick={handleRefreshData} disabled={loading}>
                  Refresh data
                </button>
                <button className="ghost" onClick={handleLogout}>
                  Logout
                </button>
              </div>
            </div>
            <p className="muted">Session is active. You can now create cafes, add reviews, and sort listings.</p>
          </section>

          <section className="card">
            <h2>Create a cafe listing</h2>
            <div className="row">
              <div>
                <label htmlFor="cafeName">Cafe name</label>
                <input id="cafeName" value={cafeName} onChange={(e) => setCafeName(e.target.value)} required />
              </div>
              <div>
                <label htmlFor="cafeAddress">Address</label>
                <input id="cafeAddress" value={cafeAddress} onChange={(e) => setCafeAddress(e.target.value)} />
              </div>
            </div>
            <div style={{ marginTop: "0.75rem" }}>
              <label htmlFor="cafeDescription">Description</label>
              <textarea
                id="cafeDescription"
                value={cafeDescription}
                onChange={(e) => setCafeDescription(e.target.value)}
                placeholder="What makes this cafe special?"
              />
            </div>
            <div className="actions">
              <button onClick={handleCreateCafe} disabled={loading || !cafeName.trim()}>
                Add cafe
              </button>
            </div>
          </section>

          <section className="card">
            <h2>Leave a review</h2>
            <div className="row three">
              <div>
                <label htmlFor="selectedCafe">Cafe</label>
                <select
                  id="selectedCafe"
                  value={selectedCafeId}
                  onChange={(e) => setSelectedCafeId(e.target.value)}
                  className="select"
                >
                  <option value="">Select a cafe</option>
                  {cafesWithMeta.map((cafe) => (
                    <option key={cafe.id} value={cafe.id}>
                      {cafe.name} (#{cafe.id})
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label htmlFor="reviewRating">Rating</label>
                <select
                  id="reviewRating"
                  value={reviewRating}
                  onChange={(e) => setReviewRating(e.target.value)}
                  className="select"
                >
                  {[5, 4, 3, 2, 1].map((score) => (
                    <option key={score} value={score}>
                      {score} / 5
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <div style={{ marginTop: "0.75rem" }}>
              <label htmlFor="reviewText">Review / notes</label>
              <textarea
                id="reviewText"
                value={reviewText}
                onChange={(e) => setReviewText(e.target.value)}
                placeholder="Write your experience..."
              />
            </div>
            <div className="actions">
              <button onClick={handleCreateReview} disabled={loading || !selectedCafeId}>
                Submit review
              </button>
            </div>
          </section>

          <section className="card">
            <div className="toolbar">
              <h2>My cafe listings</h2>
              <div>
                <label htmlFor="sortBy">Sort by</label>
                <select id="sortBy" value={sortBy} onChange={(e) => setSortBy(e.target.value)} className="select">
                  <option value="updated_desc">Recently updated</option>
                  <option value="created_desc">Newest created</option>
                  <option value="rating_desc">Highest rating</option>
                  <option value="name_asc">Name A-Z</option>
                  <option value="name_desc">Name Z-A</option>
                </select>
              </div>
            </div>
            {cafesWithMeta.length === 0 ? (
              <p className="muted">No cafes yet. Create your first listing above.</p>
            ) : (
              <div className="stack">
                {cafesWithMeta.map((cafe) => (
                  <article key={cafe.id} className="cafe-item">
                    <div className="cafe-head">
                      <div>
                        <h3>{cafe.name}</h3>
                        <p className="muted">
                          {cafe.address || "No address provided"} - avg rating:{" "}
                          {cafe.rating_count ? cafe.rating_avg.toFixed(1) : "n/a"} ({cafe.rating_count} review
                          {cafe.rating_count === 1 ? "" : "s"})
                        </p>
                      </div>
                      <div className="actions">
                        <button className="secondary" onClick={() => startCafeEdit(cafe)}>
                          Edit
                        </button>
                        <button className="danger" onClick={() => handleDeleteCafe(cafe.id)}>
                          Delete
                        </button>
                      </div>
                    </div>
                    <p>{cafe.description || "No description yet."}</p>
                  </article>
                ))}
              </div>
            )}
          </section>

          <section className="card">
            <h2>Edit cafe listing</h2>
            {!editCafeId ? (
              <p className="muted">Click Edit on a listing to load it here.</p>
            ) : (
              <>
                <div className="row">
                  <div>
                    <label htmlFor="editCafeName">Cafe name</label>
                    <input id="editCafeName" value={editCafeName} onChange={(e) => setEditCafeName(e.target.value)} />
                  </div>
                  <div>
                    <label htmlFor="editCafeAddress">Address</label>
                    <input
                      id="editCafeAddress"
                      value={editCafeAddress}
                      onChange={(e) => setEditCafeAddress(e.target.value)}
                    />
                  </div>
                </div>
                <div style={{ marginTop: "0.75rem" }}>
                  <label htmlFor="editCafeDescription">Description</label>
                  <textarea
                    id="editCafeDescription"
                    value={editCafeDescription}
                    onChange={(e) => setEditCafeDescription(e.target.value)}
                  />
                </div>
                <div className="actions">
                  <button onClick={handleUpdateCafe} disabled={loading}>
                    Save cafe changes
                  </button>
                </div>
              </>
            )}
          </section>

          <section className="card">
            <h2>My reviews</h2>
            {ratings.length === 0 ? (
              <p className="muted">No reviews yet.</p>
            ) : (
              <div className="stack">
                {ratings.map((rating) => (
                  <article key={rating.id} className="review-item">
                    <div className="toolbar">
                      <strong>
                        Cafe #{rating.cafe_listing_id} - {rating.rating}/5
                      </strong>
                      <div className="actions">
                        <button className="secondary" onClick={() => startRatingEdit(rating)}>
                          Edit
                        </button>
                        <button className="danger" onClick={() => handleDeleteRating(rating.id)}>
                          Delete
                        </button>
                      </div>
                    </div>
                    <p>{rating.review || "No review text."}</p>
                  </article>
                ))}
              </div>
            )}
          </section>

          <section className="card">
            <h2>Edit review</h2>
            {!editRatingId ? (
              <p className="muted">Click Edit on one of your reviews to modify it.</p>
            ) : (
              <>
                <div className="row three">
                  <div>
                    <label htmlFor="editRatingScore">Score</label>
                    <select
                      id="editRatingScore"
                      value={editRatingScore}
                      onChange={(e) => setEditRatingScore(e.target.value)}
                      className="select"
                    >
                      {[5, 4, 3, 2, 1].map((score) => (
                        <option key={score} value={score}>
                          {score} / 5
                        </option>
                      ))}
                    </select>
                  </div>
                </div>
                <div style={{ marginTop: "0.75rem" }}>
                  <label htmlFor="editRatingReview">Review text</label>
                  <textarea
                    id="editRatingReview"
                    value={editRatingReview}
                    onChange={(e) => setEditRatingReview(e.target.value)}
                  />
                </div>
                <div className="actions">
                  <button onClick={handleUpdateRating} disabled={loading}>
                    Save review changes
                  </button>
                </div>
              </>
            )}
          </section>
        </>
      )}

      {message ? <section className="card success">{message}</section> : null}
      {error ? <section className="card error">{error}</section> : null}
      {!isAuthed ? null : loading ? <p className="muted">Loading...</p> : null}
    </main>
  );
}
