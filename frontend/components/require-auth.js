"use client";

import Link from "next/link";
import { useAuth } from "@/components/providers/auth-provider";

export default function RequireAuth({ children }) {
  const { ready, isAuthed } = useAuth();

  if (!ready) {
    return <section className="surface empty-state">Checking your session...</section>;
  }

  if (!isAuthed) {
    return (
      <section className="surface empty-state">
        <h2>Sign in to access your personal cafe space</h2>
        <p className="muted">
          Saved places, visit tracking, and review writing are only available once you have an account.
        </p>
        <Link href="/auth" className="button">
          Go to account
        </Link>
      </section>
    );
  }

  return children;
}
