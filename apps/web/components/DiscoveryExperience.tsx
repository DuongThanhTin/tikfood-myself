"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import type { Map as MapLibreMap, Marker } from "maplibre-gl";
import { fetchDiscoveryVenues, type Venue, type VenueSearchParams } from "../lib/api";

type DiscoveryExperienceProps = {
  initialVenues: Venue[];
};

type VenueMedia = {
  cuisine: string;
  rating: string;
  reviews: string;
  badge: string;
  image: string;
  hero: string;
  videoA: string;
  videoB: string;
  viewsA: string;
  viewsB: string;
};

type IconName =
  | "bookmark"
  | "chevronLeft"
  | "chevronRight"
  | "close"
  | "explore"
  | "fire"
  | "help"
  | "location"
  | "map"
  | "play"
  | "search"
  | "share"
  | "spark"
  | "star"
  | "target"
  | "tune";

const iconGlyphs: Record<IconName, string> = {
  bookmark: "▣",
  chevronLeft: "‹",
  chevronRight: "›",
  close: "×",
  explore: "◆",
  fire: "●",
  help: "?",
  location: "⌖",
  map: "▦",
  play: "▶",
  search: "⌕",
  share: "↗",
  spark: "✦",
  star: "★",
  target: "◎",
  tune: "☰"
};

const contextChips = [
  { label: "Đang hot", value: "", kind: "trend" },
  { label: "Hẹn hò", value: "date", kind: "query" },
  { label: "Dưới 200k", value: "200000", kind: "price" },
  { label: "Ăn tối", value: "dinner", kind: "query" },
  { label: "Đi nhóm", value: "group", kind: "query" },
  { label: "Cafe chill", value: "cafe", kind: "query" }
];

const tagOptions = [
  { label: "Bánh mì", value: "banh-mi" },
  { label: "Phở", value: "pho" },
  { label: "Hẹn hò", value: "date" },
  { label: "Rooftop", value: "rooftop" },
  { label: "Street food", value: "street-food" }
];

const mediaByVenue: Record<string, VenueMedia> = {
  venue_001: {
    cuisine: "Việt Nam",
    rating: "4.8",
    reviews: "128 reviews",
    badge: "TOP 5 TRENDING",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuA9uxTUbEcCzWbUAB9x4UfK8jmlZU1kTP94dWAdj_9H7Bac5TyHS2AnUDzsNDnBViY5DonpnY7D5WaMZddbmCrhL1nJ88U3K-LtyOtCPrevPplSf7wvFjSV2-Bn-2zsoH-1SP9BjcJtDfYpbJBye7iOd69ELk0hd35ZH1h3PKzGoXFWj6-EtJZoSO2SZgpvYXU-9ugA8ioUdCU2HUHN6tnyOmV4XrlxoA4Ip8x9WR_5YuF7FMe_B3GWL7vP7woew5a1t1qWBK9A7hY",
    hero: "https://lh3.googleusercontent.com/aida-public/AB6AXuAm1vxSwpXziuE3G3Ge4t5-UkkP5G8xU7qTJzZBAxfUuQRGDAxsQb-e_-2rSeEmNPPdyOqBm0mfEvNm_BufYM_yB2FIAcpAhw88gnMnl6IhZcNS4U7xRpBuTQ5aPWK5BaZJhmpRTu_uEdRNC5uD3hKvbZV4Mi01iUeoRn50c10FvK7XbSVoFUvsEFGulxNGCu2EYMe7WEY-U7eE5D56geXsD_eKwopVurdsea8hWGiUDVcdL-a9_QEZQnQmENUNEiRA2ZU2gHn7p-U",
    videoA: "https://lh3.googleusercontent.com/aida-public/AB6AXuC06FaHZDmQ-dOX0cJbk_fOlB7J1XB67emOLcT30WikzxE44166d_4gtTWVM5ICnent7LwwD5HkmSOKq3fbiHwdtCBRZQI1R8LqZX-GVQ4pPsCg3nsFteLwS9uZ84hNGPVWA9LfJp5iTrD0C3pVcCHPQUQyIkHJ2CIOZ5pUIe8gabqJpTlAwhNTnC_x8KMWVpLgc-q0xlGZ6L-4qqhkW-DB5bXBXansRKnfdfguR6On4bjbMbCNEYYtYXUJN_4FxbZgkbLb7VF5SXw",
    videoB: "https://lh3.googleusercontent.com/aida-public/AB6AXuCan6cODRQEx7D4D9Kp-lSQRk8NS8ZNXqs5iUk8dXkJrLZ_aVRGoBTxhVsCrbEezvza5OxIILV2tkZEUEJqeR5ere6PGpiLVm5AA5GbmZr9UpXh2YIKqnnjoIqmB1psy64EUK4OdcYaH-0qI5CV_WTLEMPB7cnba3mwrlDgLraPJtSqawAjPuSRFgitm72pseNFXVqqD2m3J7pfUIFRl2FOYeDcpYttsVLQJBCieeNwOFgnx0Qz3ejNMg-SQf6VqAJrIII4c4IML98",
    viewsA: "1.2M",
    viewsB: "856K"
  },
  venue_002: {
    cuisine: "Phở",
    rating: "4.5",
    reviews: "94 reviews",
    badge: "HIDDEN GEM",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuAnO0pWvpzdbAL8SNSBFpDxRB2VqNQEcyijg79a_8IqnRmTdDcXFgQVM6kL8o8KHLB5SpL3nTTtC_ItbB0VFDpG5o9U7SATfpcXuNBQ45CCdaJvDoHHijiDSByB8Zr5dO5vY11p7fhMV7qQ8W-nsP8J4Yr_KaHfZSKnQK1FDQlm4QKM2ATjVC3mM0o5U-imGX4LjXSkoeYOpFxx_oL4V925M8NoUvg8aXZTDeW31pMAAPcXs2qa-it0ZG__0Rxp95JfaxpaF2aFeNw",
    hero: "https://lh3.googleusercontent.com/aida-public/AB6AXuA4BCVrqTVy4CSbXqFhdT3g_l3zeTjoJl8kY_PvlOfDWGZSuvcHxn-uOABWdIcO1zsVXP8kvL9rZlN97uGPpfEBWsYNFu34v69pyEIGvm-vLThxGvQ5NfiSKpQmkqbYInZxR3WHo49MGqH_7QRYQDkpwOTm5WJQmmv0IYVHe24XfrE1gTzYqTG2qNOKRvZ0AIGeU3TPr3s_KuzZd6hh7MDaMGwtX3nP1hzEsuntTrUbX5UhoQfgmFHKA68p0xWfzwyYtFS3KhdzaSY",
    videoA: "https://lh3.googleusercontent.com/aida-public/AB6AXuC06FaHZDmQ-dOX0cJbk_fOlB7J1XB67emOLcT30WikzxE44166d_4gtTWVM5ICnent7LwwD5HkmSOKq3fbiHwdtCBRZQI1R8LqZX-GVQ4pPsCg3nsFteLwS9uZ84hNGPVWA9LfJp5iTrD0C3pVcCHPQUQyIkHJ2CIOZ5pUIe8gabqJpTlAwhNTnC_x8KMWVpLgc-q0xlGZ6L-4qqhkW-DB5bXBXansRKnfdfguR6On4bjbMbCNEYYtYXUJN_4FxbZgkbLb7VF5SXw",
    videoB: "https://lh3.googleusercontent.com/aida-public/AB6AXuCan6cODRQEx7D4D9Kp-lSQRk8NS8ZNXqs5iUk8dXkJrLZ_aVRGoBTxhVsCrbEezvza5OxIILV2tkZEUEJqeR5ere6PGpiLVm5AA5GbmZr9UpXh2YIKqnnjoIqmB1psy64EUK4OdcYaH-0qI5CV_WTLEMPB7cnba3mwrlDgLraPJtSqawAjPuSRFgitm72pseNFXVqqD2m3J7pfUIFRl2FOYeDcpYttsVLQJBCieeNwOFgnx0Qz3ejNMg-SQf6VqAJrIII4c4IML98",
    viewsA: "642K",
    viewsB: "318K"
  }
};

const defaultMedia = mediaByVenue.venue_001;

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

  function applyChip(chip: (typeof contextChips)[number]) {
    if (chip.kind === "price") {
      setMaxPrice(chip.value);
      return;
    }
    if (chip.kind === "query") {
      setQuery(chip.label);
      return;
    }
    setSort("trending");
  }

  function toggleLeftPanel() {
    setLeftCollapsed((value) => !value);
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
        <button
          className="sidebarToggle"
          type="button"
          aria-label={leftCollapsed ? "Expand left panel" : "Collapse left panel"}
          onClick={toggleLeftPanel}
        >
          <Icon name={leftCollapsed ? "chevronRight" : "chevronLeft"} />
        </button>

        <div className="brandBlock">
          <div className="brandHeader">
            <span className="brandName">TikFood</span>
            {!leftCollapsed ? <span className="guestBadge">Guest View</span> : null}
          </div>
          {!leftCollapsed ? (
            <div className="locationHint">
              <Icon name="location" />
              <span>Ho Chi Minh City</span>
            </div>
          ) : null}
        </div>

        {!leftCollapsed ? (
          <>
            <nav className="sideTabs" aria-label="Primary navigation">
              <button className="sideTab active" type="button">
                <Icon name="explore" />
                Khám phá
              </button>
              <button className="sideTab" type="button" onClick={useCurrentLocation}>
                <Icon name="map" />
                Bản đồ
              </button>
            </nav>

            <section className="searchFocus">
              <h1>
                Chào mừng bạn!
                <span>Hôm nay bạn muốn đi đâu?</span>
              </h1>

              <label className="heroSearch">
                <Icon name="search" />
                <input
                  value={query}
                  onChange={(event) => setQuery(event.target.value)}
                  placeholder="VD: date dưới 500k ở Quận 1"
                />
                <Icon name="spark" className="heroSearchSpark" />
              </label>

              <div className="contextChips" aria-label="Quick discovery prompts">
                {contextChips.map((chip) => (
                  <button
                    key={chip.label}
                    className={chip.kind === "trend" ? "contextChip trend" : "contextChip"}
                    type="button"
                    onClick={() => applyChip(chip)}
                  >
                    {chip.kind === "trend" ? <Icon name="fire" /> : null}
                    {chip.label}
                  </button>
                ))}
              </div>
            </section>

            <section className="compactFilters" aria-label="Venue discovery filters">
              <div className="filterGrid">
                <label className="field">
                  <span>Quận</span>
                  <select value={district} onChange={(event) => setDistrict(event.target.value)}>
                    <option value="">Tất cả</option>
                    <option value="District 1">Quận 1</option>
                    <option value="District 3">Quận 3</option>
                  </select>
                </label>

                <label className="field">
                  <span>Giá tối đa</span>
                  <select value={maxPrice} onChange={(event) => setMaxPrice(event.target.value)}>
                    <option value="">Bất kỳ</option>
                    <option value="80000">80k</option>
                    <option value="120000">120k</option>
                    <option value="200000">200k</option>
                    <option value="300000">300k</option>
                    <option value="500000">500k</option>
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

              <div className="filterActions">
                <label className="toggle">
                  <input
                    type="checkbox"
                    checked={openNow}
                    onChange={(event) => setOpenNow(event.target.checked)}
                  />
                  <span>Đang mở cửa</span>
                </label>
                <select
                  className="sortSelect"
                  value={nearUser ? "distance" : sort}
                  onChange={(event) => setSort(event.target.value as VenueSearchParams["sort"])}
                  disabled={Boolean(nearUser)}
                  aria-label="Sort venues"
                >
                  <option value="trending">Trending</option>
                  <option value="videos">Video nhiều nhất</option>
                  <option value="price">Giá tốt</option>
                  <option value="distance">Gần nhất</option>
                </select>
              </div>

              <div className="actions">
                <button className="primaryButton" type="button" onClick={() => void runSearch()} disabled={isLoading}>
                  {isLoading ? "Đang tìm" : "Tìm kiếm"}
                </button>
                <button className="secondaryButton" type="button" onClick={clearFilters} disabled={isLoading}>
                  Reset
                </button>
              </div>

              {nearUser ? <p className="statusText">Đang lọc quanh vị trí của bạn trong 3 km.</p> : null}
              {error ? <p className="errorText">{error}</p> : null}
            </section>

            <section className="venueRail" aria-label="Recommended venues">
              <div className="sectionHeader">
                <h2>Gợi ý cho bạn</h2>
                <button type="button" onClick={() => void runSearch({ sort: "trending", limit: 20 })}>
                  Tất cả
                </button>
              </div>
              {venues.length === 0 ? <p className="emptyText">Không có địa điểm phù hợp bộ lọc.</p> : null}
              <div className="venueList">
                {venues.map((venue) => (
                  <VenueRailCard
                    key={venue.id}
                    venue={venue}
                    selected={selectedVenue?.id === venue.id}
                    onSelect={() => setSelectedVenue(venue)}
                  />
                ))}
              </div>
            </section>

            <section className="authCta">
              <p>Đăng nhập để lưu lại những địa điểm bạn yêu thích</p>
              <button type="button">Đăng nhập ngay</button>
            </section>
          </>
        ) : (
          <div className="collapsedTabs" aria-label="Collapsed navigation">
            <button className="miniTab active" type="button">
              <Icon name="explore" />
            </button>
            <button className="miniTab" type="button" onClick={useCurrentLocation}>
              <Icon name="target" />
            </button>
            <button className="miniTab" type="button" disabled>
              <Icon name="bookmark" />
            </button>
          </div>
        )}
      </aside>

      <section className="mapStage" aria-label="Restaurant map">
        <header className="mapTopbar">
          <button className="glassButton" type="button">
            <Icon name="tune" />
            Bộ lọc
          </button>
          <div className="mapStats">
            <button className="glassButton solid" type="button" onClick={useCurrentLocation}>
              Gần tôi
            </button>
            <button className="iconGlassButton" type="button" aria-label="Help">
              <Icon name="help" />
            </button>
          </div>
        </header>

        <VenueMap
          venues={venues}
          selectedVenue={selectedVenue}
          onSelectVenue={setSelectedVenue}
        />

        <div className="mapBottomAction">
          <button type="button" onClick={useCurrentLocation}>
            <Icon name="target" />
            Tìm quanh đây
          </button>
        </div>
      </section>

      {selectedVenue ? (
        <aside className="rightPanel" aria-label="Venue detail">
          <button
            className="closeButton"
            type="button"
            aria-label="Close venue detail"
            onClick={() => setSelectedVenue(null)}
          >
            <Icon name="close" />
          </button>
          <VenueDetail venue={selectedVenue} />
        </aside>
      ) : null}
    </main>
  );
}

function VenueRailCard({
  venue,
  selected,
  onSelect
}: {
  venue: Venue;
  selected: boolean;
  onSelect: () => void;
}) {
  const media = getVenueMedia(venue);
  return (
    <button
      className={selected ? "venueRailCard selected" : "venueRailCard"}
      type="button"
      onClick={onSelect}
    >
      <span className="venueThumb">
        <img alt={`${venue.name} preview`} src={media.image} />
      </span>
      <span className="venueSummary">
        <strong>{venue.name}</strong>
        <small>{media.cuisine} · {venue.district}</small>
        <span className="venueMetaRow">
          <em>{media.badge}</em>
          <span>
            <Icon name="star" />
            {media.rating}
          </span>
        </span>
      </span>
    </button>
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
        const media = getVenueMedia(venue);
        const element = document.createElement("button");
        element.className = selectedVenue?.id === venue.id ? "mapMarker selected" : "mapMarker";
        element.type = "button";
        element.setAttribute("aria-label", `Select ${venue.name}`);
        element.innerHTML = selectedVenue?.id === venue.id
          ? `<img alt="" src="${media.image}" /><span>${venue.name}</span>`
          : `<strong>${venue.trend_score}</strong>`;
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
      zoom: 14.5,
      duration: 500
    });
  }, [selectedVenue]);

  return (
    <div className="mapCanvasWrap">
      <div className="mapFallbackGrid" aria-hidden="true" />
      <div ref={mapContainerRef} className="mapCanvas" />
      <div className="mapShade" aria-hidden="true" />
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
            onClick={() => onSelectVenue(venue)}
            aria-label={`Select ${venue.name}`}
          >
            {selectedVenue?.id === venue.id ? (
              <>
                <img alt="" src={getVenueMedia(venue).image} />
                <span>{venue.name}</span>
              </>
            ) : (
              <strong>{venue.trend_score}</strong>
            )}
          </button>
        ))}
      </div>
    </div>
  );
}

function VenueDetail({ venue }: { venue: Venue }) {
  const media = getVenueMedia(venue);
  const trendPercent = Math.min(99, Math.max(1, venue.trend_score));

  return (
    <article className="detailContent">
      <div className="detailHero">
        <img alt={`${venue.name} interior`} src={media.hero} />
        <div className="detailBadges">
          <span>HOT</span>
          <span>TRENDING</span>
        </div>
      </div>

      <div className="detailBody">
        <div className="detailTitleRow">
          <div>
            <h2>{venue.name}</h2>
            <p className="detailMeta">{venue.address}, {venue.district}</p>
          </div>
          <div className="ratingBlock">
            <span>
              <Icon name="star" />
              {media.rating}
            </span>
            <small>{media.reviews}</small>
          </div>
        </div>

        <section className="trendScoreCard">
          <div>
            <span>Trend Score</span>
            <strong>{trendPercent}%</strong>
          </div>
          <div className="trendTrack">
            <span style={{ width: `${trendPercent}%` }} />
          </div>
          <p>{venue.ai_summary}</p>
        </section>

        <p className="aboutText">{venue.about}</p>

        <div className="categoryList">
          {venue.categories.map((category) => (
            <span className="category" key={category}>
              {category}
            </span>
          ))}
        </div>

        <section>
          <h3>Thịnh hành trên TikTok</h3>
          <div className="videoGrid">
            <VideoCard image={media.videoA} views={media.viewsA} label={`${venue.name} video`} />
            <VideoCard image={media.videoB} views={media.viewsB} label={`${venue.name} creator video`} />
          </div>
        </section>

        <section>
          <h3>Món đang được nhắc tới</h3>
          <div className="dishList">
            {venue.trending_dishes.map((dish) => (
              <span className="dish" key={dish}>
                {dish}
              </span>
            ))}
          </div>
        </section>

        <div className="detailActions">
          <button className="primaryButton large" type="button">Xem chi tiết & Menu</button>
          <button className="glassAction" type="button">
            <Icon name="share" />
            Chia sẻ
          </button>
          <button className="glassAction" type="button">
            <Icon name="bookmark" />
            Lưu
          </button>
        </div>
      </div>
    </article>
  );
}

function VideoCard({ image, views, label }: { image: string; views: string; label: string }) {
  return (
    <div className="videoCard">
      <img alt={label} src={image} />
      <span>
        <Icon name="play" />
        {views}
      </span>
    </div>
  );
}

function getVenueMedia(venue: Venue) {
  return mediaByVenue[venue.id] ?? defaultMedia;
}

function Icon({ name, className = "" }: { name: IconName; className?: string }) {
  return (
    <span className={["uiIcon", className].filter(Boolean).join(" ")} aria-hidden="true">
      {iconGlyphs[name]}
    </span>
  );
}
