"use client";

import { createContext, useContext, useMemo, useSyncExternalStore } from "react";

const STORAGE_KEY = "go-cafe-token";
const AuthContext = createContext(null);

const listeners = new Set();

function emitTokenChange() {
  listeners.forEach((listener) => listener());
}

function subscribe(listener) {
  listeners.add(listener);

  if (typeof window === "undefined") {
    return () => listeners.delete(listener);
  }

  const handleStorage = () => listener();
  const handleAuthChange = () => listener();
  window.addEventListener("storage", handleStorage);
  window.addEventListener("go-cafe-auth-changed", handleAuthChange);

  return () => {
    listeners.delete(listener);
    window.removeEventListener("storage", handleStorage);
    window.removeEventListener("go-cafe-auth-changed", handleAuthChange);
  };
}

function getStoredToken() {
  if (typeof window === "undefined") {
    return "";
  }

  return window.localStorage.getItem(STORAGE_KEY) || "";
}

function setStoredToken(token) {
  if (typeof window === "undefined") {
    return;
  }

  if (token) {
    window.localStorage.setItem(STORAGE_KEY, token);
  } else {
    window.localStorage.removeItem(STORAGE_KEY);
  }

  emitTokenChange();
}

export function AuthProvider({ children }) {
  const token = useSyncExternalStore(subscribe, getStoredToken, () => "");

  const value = useMemo(
    () => ({
      token,
      ready: true,
      isAuthed: Boolean(token.trim()),
      setToken: setStoredToken,
      logout() {
        setStoredToken("");
      }
    }),
    [token]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const value = useContext(AuthContext);

  if (!value) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return value;
}
