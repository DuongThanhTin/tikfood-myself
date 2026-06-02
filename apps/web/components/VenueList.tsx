import type { Venue } from "../lib/api";

export function VenueList({ venues }: { venues: Venue[] }) {
  return (
    <section className="venues">
      <h2 className="sectionTitle">Realtime discovery feed</h2>
      <div className="venueGrid">
        {venues.map((venue) => (
          <article className="venueCard" key={venue.id}>
            <div className="venueHeader">
              <div>
                <h3 className="venueName">{venue.name}</h3>
                <p className="meta">
                  {venue.address} · {venue.district}
                </p>
              </div>
              <span className="score">{venue.trend_score}</span>
            </div>
            <p>{venue.ai_summary}</p>
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
