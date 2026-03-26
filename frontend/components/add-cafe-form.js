"use client";

import { useEffect, useState } from "react";
import { searchCafeAddresses } from "@/lib/api";

const DEFAULT_IMAGE =
  "https://images.unsplash.com/photo-1509042239860-f550ce710b93?auto=format&fit=crop&w=1200&q=80";

export default function AddCafeForm({ onCreate, submitting }) {
  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [city, setCity] = useState("");
  const [neighborhood, setNeighborhood] = useState("");
  const [description, setDescription] = useState("");
  const [imageUrl, setImageUrl] = useState(DEFAULT_IMAGE);
  const [visitStatus, setVisitStatus] = useState("to_visit");
  const [latitude, setLatitude] = useState(null);
  const [longitude, setLongitude] = useState(null);
  const [suggestions, setSuggestions] = useState([]);
  const [loadingSuggestions, setLoadingSuggestions] = useState(false);

  useEffect(() => {
    const query = address.trim();

    if (query.length < 3) {
      setSuggestions([]);
      setLoadingSuggestions(false);
      return undefined;
    }

    let cancelled = false;
    const timer = window.setTimeout(async () => {
      setLoadingSuggestions(true);

      try {
        const payload = await searchCafeAddresses(query, 5);
        if (!cancelled) {
          setSuggestions(payload?.results || []);
        }
      } catch {
        if (!cancelled) {
          setSuggestions([]);
        }
      } finally {
        if (!cancelled) {
          setLoadingSuggestions(false);
        }
      }
    }, 500);

    return () => {
      cancelled = true;
      window.clearTimeout(timer);
    };
  }, [address]);

  function applySuggestion(suggestion) {
    const selectedAddress = suggestion.formatted || suggestion.address_line1 || suggestion.name || "";

    setAddress(selectedAddress);
    setCity(suggestion.city || "");
    setNeighborhood(suggestion.address_line2 || "");
    setLatitude(typeof suggestion.lat === "number" ? suggestion.lat : null);
    setLongitude(typeof suggestion.lon === "number" ? suggestion.lon : null);

    if (!name.trim() && suggestion.name) {
      setName(suggestion.name);
    }

    setSuggestions([]);
  }

  async function handleSubmit(event) {
    event.preventDefault();

    const success = await onCreate({
      name,
      address,
      city,
      neighborhood,
      description,
      image_url: imageUrl,
      visit_status: visitStatus,
      latitude,
      longitude
    });

    if (!success) {
      return;
    }

    setName("");
    setAddress("");
    setCity("");
    setNeighborhood("");
    setDescription("");
    setImageUrl(DEFAULT_IMAGE);
    setVisitStatus("to_visit");
    setLatitude(null);
    setLongitude(null);
    setSuggestions([]);
  }

  return (
    <section className="surface">
      <div className="cluster">
        <p className="eyebrow">Add a new cafe</p>
        <h2>Add a place to your personal journal</h2>
      </div>

      <form className="stack-form" onSubmit={handleSubmit}>
        <label>
          Cafe name
          <input value={name} onChange={(event) => setName(event.target.value)} required />
        </label>

        <label>
          Address
          <input value={address} onChange={(event) => setAddress(event.target.value)} />
          {loadingSuggestions ? <small className="muted">Searching addresses...</small> : null}
          {suggestions.length ? (
            <div className="suggestion-list">
              {suggestions.map((suggestion, index) => (
                <button
                  type="button"
                  key={`${suggestion.formatted || suggestion.name || "suggestion"}-${index}`}
                  className="suggestion-item"
                  onClick={() => applySuggestion(suggestion)}
                >
                  <strong>{suggestion.name || suggestion.address_line1 || "Suggested cafe"}</strong>
                  <span>{suggestion.formatted || suggestion.address_line2 || ""}</span>
                </button>
              ))}
            </div>
          ) : null}
        </label>

        <div className="two-up">
          <label>
            City
            <input value={city} onChange={(event) => setCity(event.target.value)} />
          </label>

          <label>
            Neighborhood
            <input value={neighborhood} onChange={(event) => setNeighborhood(event.target.value)} />
          </label>
        </div>

        <div className="two-up">
          <label>
            Cover image URL
            <input value={imageUrl} onChange={(event) => setImageUrl(event.target.value)} />
          </label>

          <label>
            Starting status
            <select value={visitStatus} onChange={(event) => setVisitStatus(event.target.value)}>
              <option value="to_visit">Saved for later</option>
              <option value="visited">Already visited</option>
            </select>
          </label>
        </div>

        <label>
          Why is this place worth discovering?
          <textarea
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            placeholder="Quiet weekday patio, standout pour-over, and a friendly bar team."
          />
        </label>

        <button type="submit" className="button" disabled={submitting || !name.trim()}>
          {submitting ? "Saving..." : "Add to my places"}
        </button>
      </form>
    </section>
  );
}
