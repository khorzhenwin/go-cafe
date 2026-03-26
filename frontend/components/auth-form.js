"use client";

import { useState } from "react";
import { loginUser, registerUser } from "@/lib/api";
import { useAuth } from "@/components/providers/auth-provider";

export default function AuthForm() {
  const { setToken, isAuthed } = useAuth();
  const [mode, setMode] = useState("login");
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  async function handleSubmit(event) {
    event.preventDefault();
    setLoading(true);
    setError("");
    setMessage("");

    try {
      const response =
        mode === "login"
          ? await loginUser({ email, password })
          : await registerUser({ email, name, password });

      setToken(response?.token || "");
      setMessage(mode === "login" ? "You are signed in." : "Account created and signed in.");
    } catch (submitError) {
      setError(submitError.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="surface auth-panel">
      <div className="cluster">
        <p className="eyebrow">Account</p>
        <h2>{isAuthed ? "Your session is active" : "Sign in to save places"}</h2>
        <p className="muted">
          Keep a private shortlist, mark visits, and turn discoveries into a reusable coffee journal.
        </p>
      </div>

      <div className="segmented">
        <button
          type="button"
          className={mode === "login" ? "segmented-item active" : "segmented-item"}
          onClick={() => setMode("login")}
        >
          Login
        </button>
        <button
          type="button"
          className={mode === "register" ? "segmented-item active" : "segmented-item"}
          onClick={() => setMode("register")}
        >
          Create account
        </button>
      </div>

      <form className="stack-form" onSubmit={handleSubmit}>
        <label>
          Email
          <input value={email} onChange={(event) => setEmail(event.target.value)} required />
        </label>

        {mode === "register" ? (
          <label>
            Display name
            <input value={name} onChange={(event) => setName(event.target.value)} required />
          </label>
        ) : null}

        <label>
          Password
          <input
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
        </label>

        <button type="submit" className="button" disabled={loading}>
          {loading ? "Working..." : mode === "login" ? "Login" : "Create account"}
        </button>
      </form>

      {message ? <p className="feedback success">{message}</p> : null}
      {error ? <p className="feedback error">{error}</p> : null}
    </section>
  );
}
