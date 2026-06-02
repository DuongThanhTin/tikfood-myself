export type Venue = {
  id: string;
  name: string;
  address: string;
  district: string;
  latitude: number;
  longitude: number;
  categories: string[];
  trend_score: number;
  trending_dishes: string[];
  ai_summary: string;
};

type ApiResponse<T> = {
  data: T;
  error?: string;
};

const fallbackVenues: Venue[] = [
  {
    id: "venue_001",
    name: "Banh Mi Hem",
    address: "12 Nguyen Trai",
    district: "District 1",
    latitude: 10.7712,
    longitude: 106.6899,
    categories: ["street_food", "banh_mi"],
    trend_score: 92,
    trending_dishes: ["banh mi thit nuong", "banh mi pate"],
    ai_summary: "Trending for late-night banh mi clips with strong local social proof."
  },
  {
    id: "venue_002",
    name: "Pho Bo Nguyen",
    address: "88 Le Van Sy",
    district: "District 3",
    latitude: 10.7864,
    longitude: 106.6767,
    categories: ["noodle", "pho"],
    trend_score: 87,
    trending_dishes: ["pho bo tai", "pho bo vien"],
    ai_summary: "Popular for breakfast pho videos and consistent creator mentions."
  }
];

export async function getMapVenues(): Promise<Venue[]> {
  const baseUrl = process.env.NEXT_PUBLIC_API_URL;
  if (!baseUrl) {
    return fallbackVenues;
  }

  try {
    const response = await fetch(`${baseUrl}/api/v1/map/venues`, {
      next: { revalidate: 30 }
    });
    if (!response.ok) {
      return fallbackVenues;
    }

    const body = (await response.json()) as ApiResponse<Venue[]>;
    return body.data;
  } catch {
    return fallbackVenues;
  }
}
