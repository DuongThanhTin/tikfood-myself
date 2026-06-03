import { DiscoveryExperience } from "../components/DiscoveryExperience";
import { getDiscoveryVenues } from "../lib/api";

export default async function Home() {
  const venues = await getDiscoveryVenues();

  return <DiscoveryExperience initialVenues={venues} />;
}
