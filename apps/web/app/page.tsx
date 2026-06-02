import { VenueList } from "../components/VenueList";
import { getMapVenues } from "../lib/api";

export default async function Home() {
  const venues = await getMapVenues();

  return (
    <main className="shell">
      <section className="hero">
        <div>
          <p className="eyebrow">TikFood Discovery MVP</p>
          <h1>Trending dishes near you</h1>
          <p className="lede">
            Dish-first, map-first food discovery powered by social proof,
            trend scores, and AI summaries.
          </p>
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
