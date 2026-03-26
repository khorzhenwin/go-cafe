"use client";

import Link from "next/link";
import AppShell from "@/components/app-shell";
import AuthForm from "@/components/auth-form";
import { useAuth } from "@/components/providers/auth-provider";

export default function AuthPage() {
  const { isAuthed } = useAuth();

  return (
    <AppShell
      title="Account and session"
      subtitle="Authentication is now a dedicated space, so the rest of the product can stay focused on discovery, saved places, and reviews."
      actions={
        <Link href={isAuthed ? "/my-places" : "/map"} className="button button-secondary">
          {isAuthed ? "Go to my places" : "Browse the map"}
        </Link>
      }
    >
      <section className="content-grid">
        <AuthForm />

        <section className="surface spotlight-card">
          <p className="eyebrow">What unlocks after sign in</p>
          <h2>Personal ownership without hiding discovery</h2>
          <p className="body-copy">
            Public submissions stay browsable for everyone. Signing in lets you save cafes into your own list, mark them
            visited, and add personal tasting notes.
          </p>
          <div className="cluster-sm">
            <Link href="/map" className="button">
              Explore cafes
            </Link>
            <Link href="/reviews" className="button button-ghost">
              See reviews
            </Link>
          </div>
        </section>
      </section>
    </AppShell>
  );
}
