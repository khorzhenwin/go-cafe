"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/components/providers/auth-provider";

const NAV_ITEMS = [
  { href: "/", label: "Discover" },
  { href: "/map", label: "Map" },
  { href: "/my-places", label: "My Places" },
  { href: "/reviews", label: "Reviews" }
];

export default function AppShell({ title, subtitle, actions, children }) {
  const pathname = usePathname();
  const { isAuthed, logout, ready } = useAuth();

  return (
    <div className="site-frame">
      <div className="ambient-grid" aria-hidden="true" />
      <header className="topbar surface">
        <Link href="/" className="brandmark">
          <span className="brandmark-badge">GC</span>
          <span>
            Cafe Hub
            <small>Discovery-first cafe journal</small>
          </span>
        </Link>

        <nav className="topnav" aria-label="Primary">
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={pathname === item.href ? "navlink active" : "navlink"}
            >
              {item.label}
            </Link>
          ))}
        </nav>

        <div className="topbar-actions">
          {ready && isAuthed ? (
            <>
              <Link href="/auth" className="button button-secondary">
                Account
              </Link>
              <button type="button" className="button button-ghost" onClick={logout}>
                Logout
              </button>
            </>
          ) : (
            <Link href="/auth" className="button button-secondary">
              Sign in
            </Link>
          )}
        </div>
      </header>

      <main className="page-shell">
        <section className="page-hero">
          <div>
            <p className="eyebrow">Cafe discovery platform</p>
            <h1>{title}</h1>
            {subtitle ? <p className="lead">{subtitle}</p> : null}
          </div>
          {actions ? <div className="page-hero-actions">{actions}</div> : null}
        </section>

        {children}
      </main>
    </div>
  );
}
