"use client";

import { useEffect, useMemo, useState } from "react";
import {
  createCafeRating,
  createMyCafe,
  deleteCafe,
  deleteRating,
  listMyCafes,
  listMyRatings,
  loginUser,
  registerUser,
  updateCafe
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
  const [pendingStatusById, setPendingStatusById] = useState({});

  const [cafeName, setCafeName] = useState("");
  const [cafeAddress, setCafeAddress] = useState("");
  const [cafeDescription, setCafeDescription] = useState("");
  const [newCafeStatus, setNewCafeStatus] = useState("to_visit");

  const [selectedCafeId, setSelectedCafeId] = useState("");
  const [reviewRating, setReviewRating] = useState("5");
  const [reviewText, setReviewText] = useState("");

  const [sortBy, setSortBy] = useState("updated_desc");
  const [activeStep, setActiveStep] = useState(1);

  const isAuthed = useMemo(() => token.trim().length > 0, [token]);
  const toVisitCafes = useMemo(() => cafes.filter((c) => c.visit_status === "to_visit"), [cafes]);
  const visitedCafes = useMemo(() => cafes.filter((c) => c.visit_status === "visited"), [cafes]);

  useEffect(() => {
    const saved = window.localStorage.getItem("go-cafe-token");
    if (saved) setToken(saved);
  }, []);

  useEffect(() => {
    if (token) window.localStorage.setItem("go-cafe-token", token);
    else window.localStorage.removeItem("go-cafe-token");
  }, [token]);

  useEffect(() => {
    if (isAuthed) {
      handleRefreshData();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthed]);

  useEffect(() => {
    if (isAuthed) {
      handleRefreshData();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sortBy]);

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
      const [myCafes, myRatings] = await Promise.all([listMyCafes(token, { sort: sortBy }), listMyRatings(token)]);
      const cafesData = myCafes || [];
      setCafes(cafesData);
      setRatings(myRatings || []);
      setSelectedCafeId((prev) => prev || String(cafesData[0]?.id || ""));
      setPendingStatusById(
        cafesData.reduce((acc, cafe) => {
          acc[cafe.id] = cafe.visit_status || "to_visit";
          return acc;
        }, {})
      );
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
      await createMyCafe(token, {
        name: cafeName,
        address: cafeAddress,
        description: cafeDescription,
        visit_status: newCafeStatus
      });
      setCafeName("");
      setCafeAddress("");
      setCafeDescription("");
      setNewCafeStatus("to_visit");
      setMessage("Cafe listing created.");
      await handleRefreshData();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleUpdateCafeStatus(cafe) {
    setError("");
    setMessage("");
    setLoading(true);
    try {
      await updateCafe(token, cafe.id, {
        name: cafe.name,
        address: cafe.address,
        description: cafe.description,
        visit_status: pendingStatusById[cafe.id] || cafe.visit_status
      });
      setMessage("Cafe status updated.");
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
    setPendingStatusById({});
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
      <div className="splash-screen" aria-hidden="true">
        <div className="vector-grid" />
        <div className="vector-orb orb-a" />
        <div className="vector-orb orb-b" />
        <div className="vector-orb orb-c" />
        <div className="vector-line line-a" />
        <div className="vector-line line-b" />
        <div className="vector-line line-c" />
      </div>
      <header className="hero">
        <div>
          <h1>Cafe Hub</h1>
          <p className="muted">
            A simple cafe journey: <strong>add cafes</strong>, <strong>mark status</strong>, then{" "}
            <strong>rate visited cafes</strong>.
          </p>
        </div>
      </header>

      {!isAuthed ? (
        <section className="flow-section">
          <form onSubmit={handleAuthSubmit}>
            <div className="stack-fields">
              <label htmlFor="email">
                Email
                <input id="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
              </label>
              {authMode === "signup" ? (
                <label htmlFor="name">
                  Display name
                  <input id="name" value={name} onChange={(e) => setName(e.target.value)} required />
                </label>
              ) : null}
              <label htmlFor="password">
                Password
                <input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </label>
            </div>
            <div className="actions">
              <button type="submit" disabled={loading}>
                {loading ? "Please wait..." : authMode === "signup" ? "Create account" : "Login"}
              </button>
            </div>
            <div className="auth-tabs auth-tabs-bottom">
              <button className={authMode === "signup" ? "tab active" : "tab"} onClick={() => setAuthMode("signup")}>
                Sign up
              </button>
              <button className={authMode === "login" ? "tab active" : "tab"} onClick={() => setAuthMode("login")}>
                Login
              </button>
            </div>
          </form>
        </section>
      ) : (
        <>
          <section className="flow-section">
            <div className="toolbar">
              <h2>Journey</h2>
              <div className="actions">
                <button className="secondary" onClick={handleRefreshData} disabled={loading}>
                  Refresh data
                </button>
                <button className="ghost" onClick={handleLogout}>
                  Logout
                </button>
              </div>
            </div>
            <p className="muted">Follow steps in order. Reviews are available only after a cafe is marked visited.</p>
            <div className="step-tabs" role="tablist" aria-label="Journey steps">
              <button
                className={activeStep === 1 ? "tab active" : "tab"}
                onClick={() => setActiveStep(1)}
                role="tab"
                aria-selected={activeStep === 1}
              >
                Step 1 - Add cafe
              </button>
              <button
                className={activeStep === 2 ? "tab active" : "tab"}
                onClick={() => setActiveStep(2)}
                role="tab"
                aria-selected={activeStep === 2}
              >
                Step 2 - Update status
              </button>
              <button
                className={activeStep === 3 ? "tab active" : "tab"}
                onClick={() => setActiveStep(3)}
                role="tab"
                aria-selected={activeStep === 3}
              >
                Step 3 - Rate visited
              </button>
            </div>
          </section>

          {activeStep === 1 ? (
            <section className="flow-section">
              <h2>Step 1 - Add your cafe</h2>
              <div className="stack-fields">
                <label>
                  Cafe name
                  <input value={cafeName} onChange={(e) => setCafeName(e.target.value)} required />
                </label>
                <label>
                  Address
                  <input value={cafeAddress} onChange={(e) => setCafeAddress(e.target.value)} />
                </label>
                <label>
                  Initial status
                  <select value={newCafeStatus} onChange={(e) => setNewCafeStatus(e.target.value)} className="select">
                    <option value="to_visit">To Visit</option>
                    <option value="visited">Visited</option>
                  </select>
                </label>
                <label>
                  Description
                  <textarea
                    value={cafeDescription}
                    onChange={(e) => setCafeDescription(e.target.value)}
                    placeholder="Why this cafe matters to you..."
                  />
                </label>
              </div>
              <div className="actions">
                <button onClick={handleCreateCafe} disabled={loading || !cafeName.trim()}>
                  Add cafe
                </button>
                <button className="secondary" onClick={() => setActiveStep(2)} disabled={loading}>
                  Continue to Step 2
                </button>
              </div>
            </section>
          ) : null}

          {activeStep === 2 ? (
            <section className="flow-section">
              <div className="toolbar">
                <h2>Step 2 - Update journey status</h2>
                <label>
                  Sort
                  <select value={sortBy} onChange={(e) => setSortBy(e.target.value)} className="select">
                    <option value="updated_desc">Recently updated</option>
                    <option value="created_desc">Newest created</option>
                    <option value="name_asc">Name A-Z</option>
                    <option value="name_desc">Name Z-A</option>
                    <option value="status_asc">Status A-Z</option>
                    <option value="status_desc">Status Z-A</option>
                  </select>
                </label>
              </div>
              <p className="muted">Mark cafes as visited before adding ratings.</p>
              {cafes.length === 0 ? (
                <p className="muted">No cafes yet. Complete Step 1 first.</p>
              ) : (
                <div className="status-layout">
                  <div>
                    <h3>To Visit</h3>
                    {toVisitCafes.length === 0 ? (
                      <p className="muted">None</p>
                    ) : (
                      toVisitCafes.map((cafe) => (
                        <div key={cafe.id} className="list-row">
                          <div>
                            <strong>{cafe.name}</strong>
                            <p className="muted">{cafe.address || "No address"}</p>
                          </div>
                          <div className="row-actions">
                            <select
                              className="select"
                              value={pendingStatusById[cafe.id] || cafe.visit_status}
                              onChange={(e) =>
                                setPendingStatusById((prev) => ({ ...prev, [cafe.id]: e.target.value }))
                              }
                            >
                              <option value="to_visit">To Visit</option>
                              <option value="visited">Visited</option>
                            </select>
                            <button
                              className="secondary"
                              onClick={() => handleUpdateCafeStatus(cafe)}
                              disabled={loading}
                            >
                              Save
                            </button>
                            <button className="ghost" onClick={() => handleDeleteCafe(cafe.id)} disabled={loading}>
                              Remove
                            </button>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                  <div>
                    <h3>Visited</h3>
                    {visitedCafes.length === 0 ? (
                      <p className="muted">None yet</p>
                    ) : (
                      visitedCafes.map((cafe) => (
                        <div key={cafe.id} className="list-row">
                          <div>
                            <strong>{cafe.name}</strong>
                            <p className="muted">{cafe.address || "No address"}</p>
                          </div>
                          <div className="row-actions">
                            <select
                              className="select"
                              value={pendingStatusById[cafe.id] || cafe.visit_status}
                              onChange={(e) =>
                                setPendingStatusById((prev) => ({ ...prev, [cafe.id]: e.target.value }))
                              }
                            >
                              <option value="to_visit">To Visit</option>
                              <option value="visited">Visited</option>
                            </select>
                            <button
                              className="secondary"
                              onClick={() => handleUpdateCafeStatus(cafe)}
                              disabled={loading}
                            >
                              Save
                            </button>
                            <button className="ghost" onClick={() => handleDeleteCafe(cafe.id)} disabled={loading}>
                              Remove
                            </button>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                </div>
              )}
              <div className="actions">
                <button className="ghost" onClick={() => setActiveStep(1)} disabled={loading}>
                  Back to Step 1
                </button>
                <button className="secondary" onClick={() => setActiveStep(3)} disabled={loading}>
                  Continue to Step 3
                </button>
              </div>
            </section>
          ) : null}

          {activeStep === 3 ? (
            <section className="flow-section">
              <h2>Step 3 - Rate visited cafes</h2>
              <div className="stack-fields">
                <label>
                  Cafe
                  <select value={selectedCafeId} onChange={(e) => setSelectedCafeId(e.target.value)} className="select">
                    <option value="">Select visited cafe</option>
                    {visitedCafes.map((cafe) => (
                      <option key={cafe.id} value={cafe.id}>
                        {cafe.name} (#{cafe.id})
                      </option>
                    ))}
                  </select>
                </label>
                <label>
                  Rating
                  <select value={reviewRating} onChange={(e) => setReviewRating(e.target.value)} className="select">
                    {[5, 4, 3, 2, 1].map((score) => (
                      <option key={score} value={score}>
                        {score}/5
                      </option>
                    ))}
                  </select>
                </label>
                <label>
                  Review
                  <textarea
                    value={reviewText}
                    onChange={(e) => setReviewText(e.target.value)}
                    placeholder="How was your visit?"
                  />
                </label>
              </div>
              <div className="actions">
                <button onClick={handleCreateReview} disabled={loading || !selectedCafeId || visitedCafes.length === 0}>
                  Submit review
                </button>
                <button className="ghost" onClick={() => setActiveStep(2)} disabled={loading}>
                  Back to Step 2
                </button>
              </div>
              {visitedCafes.length === 0 ? (
                <p className="muted">No visited cafes yet. Complete Step 2 first.</p>
              ) : null}
            </section>
          ) : null}

          <section className="flow-section">
            <h2>Recent reviews</h2>
            {ratings.length === 0 ? (
              <p className="muted">No reviews yet.</p>
            ) : (
              <div className="review-list">
                {ratings.map((rating) => (
                  <article key={rating.id} className="review-line">
                    <strong>
                      Cafe #{rating.cafe_listing_id} - {rating.rating}/5
                    </strong>
                    <p>{rating.review || "No review text."}</p>
                    <button className="ghost" onClick={() => handleDeleteRating(rating.id)} disabled={loading}>
                      Delete
                    </button>
                  </article>
                ))}
              </div>
            )}
          </section>
        </>
      )}

      {message ? <section className="flow-section success">{message}</section> : null}
      {error ? <section className="flow-section error">{error}</section> : null}
      {!isAuthed ? null : loading ? <p className="muted">Loading...</p> : null}
    </main>
  );
}
