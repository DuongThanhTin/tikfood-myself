import type { Venue } from "../lib/api";

export function VenueList({ venues }: { venues: Venue[] }) {
  return (
    <section className="venues">
      <div className="sectionHeader">
        <h2 className="sectionTitle">Realtime discovery feed</h2>
        <span className="resultCount">{venues.length} venues</span>
      </div>
      {venues.length === 0 ? (
        <div className="emptyState">No venues match the current filters.</div>
      ) : null}
      <div className="venueGrid">
        {venues.map((venue) => (
          <article className="venueCard" key={venue.id}>
            <div className="venueHeader">
              <div>
                <h3 className="venueName">{venue.name}</h3>
                <p className="meta">{venue.short_description}</p>
                <p className="meta">
                  {venue.address} · {venue.district}
                </p>
                <p className="meta">
                  {formatPrice(venue.avg_price_min_vnd)} - {formatPrice(venue.avg_price_max_vnd)} · {venue.social_video_count} videos
                  {venue.distance_meters !== undefined ? ` · ${formatDistance(venue.distance_meters)}` : ""}
                </p>
              </div>
              <span className="score">{venue.trend_score}</span>
            </div>
            <p>{venue.about}</p>
            <p>{venue.ai_summary}</p>
            <div className="categoryList">
              {venue.categories.map((category) => (
                <span className="category" key={category}>
                  {category}
                </span>
              ))}
            </div>
            <div className="dishList">
              {venue.trending_dishes.map((dish) => (
                <span className="dish" key={dish}>
                  {dish}
                </span>
              ))}
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}

function formatPrice(value: number) {
  if (value >= 1000) {
    return `${Math.round(value / 1000)}k`;
  }
  return String(value);
}

function formatDistance(value: number) {
  if (value >= 1000) {
    return `${(value / 1000).toFixed(1)} km`;
  }
  return `${Math.round(value)} m`;
}
