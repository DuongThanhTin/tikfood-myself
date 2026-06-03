"use client";

import { useMemo, useState, useTransition } from "react";
import { VenueList } from "./VenueList";
import { fetchDiscoveryVenues, type Venue, type VenueSearchParams } from "../lib/api";

type DiscoveryExperienceProps = {
  initialVenues: Venue[];
};

const tagOptions = [
  { label: "Banh mi", value: "banh-mi" },
  { label: "Pho", value: "pho" },
  { label: "Breakfast", value: "breakfast" },
  { label: "Late night", value: "late-night" },
  { label: "Street food", value: "street-food" }
];

export function DiscoveryExperience({ initialVenues }: DiscoveryExperienceProps) {
  const [venues, setVenues] = useState(initialVenues);
  const [query, setQuery] = useState("");
  const [district, setDistrict] = useState("");
  const [maxPrice, setMaxPrice] = useState("120000");
  const [activeTag, setActiveTag] = useState("");
  const [sort, setSort] = useState<VenueSearchParams["sort"]>("trending");
  const [openNow, setOpenNow] = useState(false);
  const [error, setError] = useState("");
  const [isPending, startTransition] = useTransition();

  const params = useMemo<VenueSearchParams>(() => ({
    q: query,
    district,
    tags: activeTag,
    max_price_vnd: Number(maxPrice) || undefined,
    open_now: openNow,
    sort,
    limit: 20
  }), [activeTag, district, maxPrice, openNow, query, sort]);

  function runSearch(nextParams = params) {
    setError("");
    startTransition(() => {
      void (async () => {
        try {
          const nextVenues = await fetchDiscoveryVenues(nextParams);
          setVenues(nextVenues);
        } catch (searchError) {
          setError(searchError instanceof Error ? searchError.message : "Failed to load venues.");
        }
      })();
    });
  }

  function clearFilters() {
    const nextParams: VenueSearchParams = {
      sort: "trending",
      limit: 20
    };
    setQuery("");
    setDistrict("");
    setMaxPrice("");
    setActiveTag("");
    setSort("trending");
    setOpenNow(false);
    runSearch(nextParams);
  }

  return (
    <main className="shell">
      <section className="workspace">
        <div className="searchPanel">
          <p className="eyebrow">TikFood Discovery MVP</p>
          <h1>Find trending food by dish, place, price, and social proof.</h1>

          <div className="filters" aria-label="Venue discovery filters">
            <label className="field">
              <span>Search</span>
              <input
                value={query}
                onChange={(event) => setQuery(event.target.value)}
                placeholder="pho, banh mi, date, rooftop"
              />
            </label>

            <div className="filterRow">
              <label className="field">
                <span>District</span>
                <select value={district} onChange={(event) => setDistrict(event.target.value)}>
                  <option value="">All</option>
                  <option value="District 1">District 1</option>
                  <option value="District 3">District 3</option>
                </select>
              </label>

              <label className="field">
                <span>Max price</span>
                <select value={maxPrice} onChange={(event) => setMaxPrice(event.target.value)}>
                  <option value="">Any</option>
                  <option value="80000">80k</option>
                  <option value="120000">120k</option>
                  <option value="300000">300k</option>
                  <option value="500000">500k</option>
                </select>
              </label>

              <label className="field">
                <span>Sort</span>
                <select value={sort} onChange={(event) => setSort(event.target.value as VenueSearchParams["sort"])}>
                  <option value="trending">Trending</option>
                  <option value="videos">Videos</option>
                  <option value="price">Price</option>
                </select>
              </label>
            </div>

            <div className="tagBar" aria-label="Tags">
              {tagOptions.map((tag) => (
                <button
                  key={tag.value}
                  type="button"
                  className={activeTag === tag.value ? "tag active" : "tag"}
                  onClick={() => setActiveTag(activeTag === tag.value ? "" : tag.value)}
                >
                  {tag.label}
                </button>
              ))}
            </div>

            <label className="toggle">
              <input
                type="checkbox"
                checked={openNow}
                onChange={(event) => setOpenNow(event.target.checked)}
              />
              <span>Open now</span>
            </label>

            <div className="actions">
              <button className="primaryButton" type="button" onClick={() => runSearch()} disabled={isPending}>
                {isPending ? "Searching" : "Search"}
              </button>
              <button className="secondaryButton" type="button" onClick={clearFilters} disabled={isPending}>
                Reset
              </button>
            </div>

            {error ? <p className="errorText">{error}</p> : null}
          </div>
        </div>

        <div className="mapPanel" aria-label="Discovery map preview">
          {venues.map((venue) => (
            <span
              key={venue.id}
              className="marker"
              style={{
                left: `${Math.max(8, Math.min(88, (venue.longitude - 106.65) * 420))}%`,
                top: `${Math.max(10, Math.min(82, 90 - (venue.latitude - 10.74) * 900))}%`
              }}
              title={venue.name}
            >
              {venue.trend_score}
            </span>
          ))}
        </div>
      </section>

      <VenueList venues={venues} />
    </main>
  );
}
