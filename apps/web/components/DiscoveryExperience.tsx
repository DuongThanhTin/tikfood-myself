"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import type { Map as MapLibreMap, Marker } from "maplibre-gl";
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
  const [selectedVenue, setSelectedVenue] = useState<Venue | null>(initialVenues[0] ?? null);
  const [leftCollapsed, setLeftCollapsed] = useState(false);
  const [query, setQuery] = useState("");
  const [district, setDistrict] = useState("");
  const [maxPrice, setMaxPrice] = useState("120000");
  const [activeTag, setActiveTag] = useState("");
  const [sort, setSort] = useState<VenueSearchParams["sort"]>("trending");
  const [openNow, setOpenNow] = useState(false);
  const [nearUser, setNearUser] = useState<{ lat: number; lng: number } | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const params = useMemo<VenueSearchParams>(() => ({
    q: query,
    district,
    tags: activeTag,
    lat: nearUser?.lat,
    lng: nearUser?.lng,
    radius_m: nearUser ? 3000 : undefined,
    max_price_vnd: Number(maxPrice) || undefined,
    open_now: openNow,
    sort: nearUser ? "distance" : sort,
    limit: 20
  }), [activeTag, district, maxPrice, nearUser, openNow, query, sort]);

  async function runSearch(nextParams = params) {
    setError("");
    setIsLoading(true);
    try {
      const nextVenues = await fetchDiscoveryVenues(nextParams);
      setVenues(nextVenues);
      setSelectedVenue(nextVenues[0] ?? null);
    } catch (searchError) {
      setError(searchError instanceof Error ? searchError.message : "Failed to load venues.");
    } finally {
      setIsLoading(false);
    }
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
    setNearUser(null);
    void runSearch(nextParams);
  }

  function useCurrentLocation() {
    if (!navigator.geolocation) {
      setError("Your browser does not support location sharing.");
      return;
    }

    setError("");
    setIsLoading(true);
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const nextLocation = {
          lat: position.coords.latitude,
          lng: position.coords.longitude
        };
        setNearUser(nextLocation);
        void runSearch({
          ...params,
          lat: nextLocation.lat,
          lng: nextLocation.lng,
          radius_m: 3000,
          sort: "distance"
        });
      },
      () => {
        setError("Location permission was not granted.");
        setIsLoading(false);
      },
      { enableHighAccuracy: true, timeout: 8000 }
    );
  }

  return (
    <main
      className={[
        "appShell",
        leftCollapsed ? "leftCollapsed" : "",
        selectedVenue ? "detailOpen" : ""
      ].filter(Boolean).join(" ")}
    >
      <aside className="leftPanel" aria-label="Discovery controls">
        <div className="brandRow">
          <div className="brandMark">T</div>
          {!leftCollapsed ? (
            <div>
              <p className="brandName">TikFood</p>
              <p className="brandMode">Guest discovery</p>
            </div>
          ) : null}
          <button
            className="iconButton"
            type="button"
            aria-label={leftCollapsed ? "Expand left panel" : "Collapse left panel"}
            onMouseDown={() => setLeftCollapsed(!leftCollapsed)}
            onClick={() => setLeftCollapsed(!leftCollapsed)}
          >
            {leftCollapsed ? ">" : "<"}
          </button>
        </div>

        {!leftCollapsed ? (
          <>
            <nav className="sideTabs" aria-label="Primary navigation">
              <button className="sideTab active" type="button">Explore</button>
              <button className="sideTab" type="button" onClick={useCurrentLocation}>Near me</button>
              <button className="sideTab muted" type="button" disabled>Saved</button>
              <button className="sideTab muted" type="button" disabled>Collections</button>
              <button className="sideTab muted" type="button" disabled>Profile</button>
            </nav>

            <section className="controlBlock">
              <div>
                <p className="eyebrow">Homepage</p>
                <h1>Discover where to eat now.</h1>
                <p className="panelCopy">
                  Trending restaurants and dishes stay visible before login. Location sharing narrows results nearby.
                </p>
              </div>

              <label className="field">
                <span>Search</span>
                <input
                  value={query}
                  onChange={(event) => setQuery(event.target.value)}
                  placeholder="pho, banh mi, date, rooftop"
                />
              </label>

              <div className="filterGrid">
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
              </div>

              <label className="field">
                <span>Sort</span>
                <select
                  value={nearUser ? "distance" : sort}
                  onChange={(event) => setSort(event.target.value as VenueSearchParams["sort"])}
                  disabled={Boolean(nearUser)}
                >
                  <option value="trending">Trending</option>
                  <option value="videos">Videos</option>
                  <option value="price">Price</option>
                  <option value="distance">Distance</option>
                </select>
              </label>

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
                <button className="primaryButton" type="button" onClick={() => void runSearch()} disabled={isLoading}>
                  {isLoading ? "Searching" : "Search"}
                </button>
                <button className="secondaryButton" type="button" onClick={clearFilters} disabled={isLoading}>
                  Reset
                </button>
              </div>

              {nearUser ? <p className="statusText">Location filter active within 3 km.</p> : null}
              {error ? <p className="errorText">{error}</p> : null}
            </section>

            <section className="venueRail" aria-label="Venue results">
              <div className="sectionHeader">
                <h2>Top places</h2>
                <span>{venues.length}</span>
              </div>
              {venues.length === 0 ? <p className="emptyText">No venues match the filters.</p> : null}
              {venues.map((venue) => (
                <button
                  className={selectedVenue?.id === venue.id ? "venueRailCard selected" : "venueRailCard"}
                  key={venue.id}
                  type="button"
                  onMouseDown={() => setSelectedVenue(venue)}
                  onClick={() => setSelectedVenue(venue)}
                >
                  <span>
                    <strong>{venue.name}</strong>
                    <small>{venue.district} · {formatPrice(venue.avg_price_max_vnd)}</small>
                  </span>
                  <em>{venue.trend_score}</em>
                </button>
              ))}
            </section>
          </>
        ) : (
          <div className="collapsedTabs" aria-label="Collapsed navigation">
            <button className="miniTab active" type="button">E</button>
            <button className="miniTab" type="button" onClick={useCurrentLocation}>N</button>
            <button className="miniTab" type="button" disabled>S</button>
            <button className="miniTab" type="button" disabled>C</button>
            <button className="miniTab" type="button" disabled>P</button>
          </div>
        )}
      </aside>

      <section className="mapStage" aria-label="Restaurant map">
        <div className="mapTopbar">
          <div>
            <p className="eyebrow">Live map</p>
            <h2>Trending near Ho Chi Minh City</h2>
          </div>
          <div className="mapStats">
            <span>{venues.length} places</span>
            <span>{selectedVenue ? selectedVenue.name : "Select a marker"}</span>
          </div>
        </div>
        <VenueMap
          venues={venues}
          selectedVenue={selectedVenue}
          onSelectVenue={setSelectedVenue}
        />
      </section>

      {selectedVenue ? (
        <aside className="rightPanel" aria-label="Venue detail">
          <button
            className="closeButton"
            type="button"
            onMouseDown={() => setSelectedVenue(null)}
            onClick={() => setSelectedVenue(null)}
          >
            Close
          </button>
          <VenueDetail venue={selectedVenue} />
        </aside>
      ) : null}
    </main>
  );
}

function VenueMap({
  venues,
  selectedVenue,
  onSelectVenue
}: {
  venues: Venue[];
  selectedVenue: Venue | null;
  onSelectVenue: (venue: Venue) => void;
}) {
  const mapContainerRef = useRef<HTMLDivElement | null>(null);
  const mapRef = useRef<MapLibreMap | null>(null);
  const markersRef = useRef<Marker[]>([]);

  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) {
      return;
    }

    let cancelled = false;
    void import("maplibre-gl").then((maplibregl) => {
      if (cancelled || !mapContainerRef.current) {
        return;
      }

      const map = new maplibregl.Map({
        container: mapContainerRef.current,
        center: [106.683, 10.778],
        zoom: 13,
        attributionControl: false,
        style: {
          version: 8,
          sources: {
            osm: {
              type: "raster",
              tiles: ["https://tile.openstreetmap.org/{z}/{x}/{y}.png"],
              tileSize: 256,
              attribution: "OpenStreetMap"
            }
          },
          layers: [
            {
              id: "osm",
              type: "raster",
              source: "osm"
            }
          ]
        }
      });

      map.addControl(new maplibregl.NavigationControl({ visualizePitch: true }), "bottom-right");
      mapRef.current = map;
    });

    return () => {
      cancelled = true;
      markersRef.current.forEach((marker) => marker.remove());
      markersRef.current = [];
      mapRef.current?.remove();
      mapRef.current = null;
    };
  }, []);

  useEffect(() => {
    const map = mapRef.current;
    if (!map) {
      return;
    }

    markersRef.current.forEach((marker) => marker.remove());
    markersRef.current = [];

    void import("maplibre-gl").then((maplibregl) => {
      venues.forEach((venue) => {
        const element = document.createElement("button");
        element.className = selectedVenue?.id === venue.id ? "mapMarker selected" : "mapMarker";
        element.type = "button";
        element.textContent = String(venue.trend_score);
        element.setAttribute("aria-label", `Select ${venue.name}`);
        element.addEventListener("click", () => onSelectVenue(venue));

        const marker = new maplibregl.Marker({ element })
          .setLngLat([venue.longitude, venue.latitude])
          .addTo(map);
        markersRef.current.push(marker);
      });
    });
  }, [onSelectVenue, selectedVenue, venues]);

  useEffect(() => {
    const map = mapRef.current;
    if (!map || !selectedVenue) {
      return;
    }
    map.flyTo({
      center: [selectedVenue.longitude, selectedVenue.latitude],
      zoom: 15,
      duration: 500
    });
  }, [selectedVenue]);

  return (
    <div className="mapCanvasWrap">
      <div className="mapFallbackGrid" aria-hidden="true" />
      <div ref={mapContainerRef} className="mapCanvas" />
      <div className="fallbackMarkerLayer" aria-label="Restaurant markers">
        {venues.map((venue) => (
          <button
            key={venue.id}
            className={selectedVenue?.id === venue.id ? "fallbackMarker selected" : "fallbackMarker"}
            type="button"
            style={{
              left: `${Math.max(12, Math.min(88, (venue.longitude - 106.65) * 420))}%`,
              top: `${Math.max(14, Math.min(82, 90 - (venue.latitude - 10.74) * 900))}%`
            }}
            onMouseDown={() => onSelectVenue(venue)}
            onClick={() => onSelectVenue(venue)}
            aria-label={`Select ${venue.name}`}
          >
            {venue.trend_score}
          </button>
        ))}
      </div>
    </div>
  );
}

function VenueDetail({ venue }: { venue: Venue }) {
  return (
    <article className="detailContent">
      <div>
        <p className="eyebrow">Restaurant detail</p>
        <h2>{venue.name}</h2>
        <p className="detailMeta">{venue.address} · {venue.district}</p>
      </div>

      <div className="detailScoreGrid">
        <span>
          <strong>{venue.trend_score}</strong>
          Trend
        </span>
        <span>
          <strong>{venue.social_video_count}</strong>
          Videos
        </span>
        <span>
          <strong>{formatPrice(venue.avg_price_max_vnd)}</strong>
          Max price
        </span>
      </div>

      <p>{venue.about}</p>
      <p className="aiSummary">{venue.ai_summary}</p>

      <div className="categoryList">
        {venue.categories.map((category) => (
          <span className="category" key={category}>
            {category}
          </span>
        ))}
      </div>

      <div>
        <h3>Trending dishes</h3>
        <div className="dishList">
          {venue.trending_dishes.map((dish) => (
            <span className="dish" key={dish}>
              {dish}
            </span>
          ))}
        </div>
      </div>

      <div className="authHint">
        <strong>Login later</strong>
        <p>Saved, collections, and profile tabs will appear for authenticated users.</p>
      </div>
    </article>
  );
}

function formatPrice(value: number) {
  if (value >= 1000) {
    return `${Math.round(value / 1000)}k`;
  }
  return String(value);
}
