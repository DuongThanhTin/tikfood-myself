export type Venue = {
  id: string;
  name: string;
  slug: string;
  short_description: string;
  about: string;
  address: string;
  city: string;
  district: string;
  latitude: number;
  longitude: number;
  categories: string[];
  price_level: number;
  avg_price_min_vnd: number;
  avg_price_max_vnd: number;
  currency: string;
  social_video_count: number;
  trend_score: number;
  trending_dishes: string[];
  ai_summary: string;
  distance_meters?: number;
};

type ApiResponse<T> = {
  data: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
};

export type VenueSearchParams = {
  q?: string;
  district?: string;
  dish?: string;
  tags?: string;
  lat?: number;
  lng?: number;
  radius_m?: number;
  max_price_vnd?: number;
  open_now?: boolean;
  sort?: "trending" | "videos" | "distance" | "price";
  limit?: number;
};

export const fallbackVenues: Venue[] = [
  {
    id: "venue_001",
    name: "Banh Mi Hem",
    slug: "banh-mi-hem-nguyen-trai-district-1",
    short_description: "Late-night banh mi spot trending on social video.",
    about: "A compact street-food venue known for grilled pork banh mi, quick service, and strong late-night local buzz near Nguyen Trai.",
    address: "12 Nguyen Trai",
    city: "Ho Chi Minh City",
    district: "District 1",
    latitude: 10.7712,
    longitude: 106.6899,
    categories: ["banh-mi", "late-night", "street-food", "vietnamese"],
    price_level: 1,
    avg_price_min_vnd: 30000,
    avg_price_max_vnd: 80000,
    currency: "VND",
    social_video_count: 42,
    trend_score: 92,
    trending_dishes: ["banh mi thit nuong", "banh mi pate"],
    ai_summary: "Trending for late-night banh mi clips with strong local social proof."
  },
  {
    id: "venue_002",
    name: "Pho Bo Nguyen",
    slug: "pho-bo-nguyen-le-van-sy-district-3",
    short_description: "Breakfast pho shop with consistent creator mentions.",
    about: "A neighborhood pho venue known for clear broth, beef toppings, and steady breakfast traffic from local regulars and food creators.",
    address: "88 Le Van Sy",
    city: "Ho Chi Minh City",
    district: "District 3",
    latitude: 10.7864,
    longitude: 106.6767,
    categories: ["breakfast", "pho", "vietnamese"],
    price_level: 1,
    avg_price_min_vnd: 50000,
    avg_price_max_vnd: 120000,
    currency: "VND",
    social_video_count: 35,
    trend_score: 87,
    trending_dishes: ["pho bo tai", "pho bo vien"],
    ai_summary: "Popular for breakfast pho videos and consistent creator mentions."
  }
];

export async function getDiscoveryVenues(params: VenueSearchParams = {}): Promise<Venue[]> {
  const baseUrl = process.env.NEXT_PUBLIC_API_URL;
  if (!baseUrl) {
    return fallbackVenues;
  }

  try {
    const query = toSearchParams(params);
    const response = await fetch(`${baseUrl}/api/v1/discovery/venues${query}`, {
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

export function getClientApiBaseUrl() {
  return process.env.NEXT_PUBLIC_API_URL ?? "";
}

export async function fetchDiscoveryVenues(params: VenueSearchParams): Promise<Venue[]> {
  const baseUrl = getClientApiBaseUrl();
  if (!baseUrl) {
    return filterFallbackVenues(params);
  }

  const response = await fetch(`${baseUrl}/api/v1/discovery/venues${toSearchParams(params)}`);
  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as ApiResponse<null> | null;
    throw new Error(body?.error?.message ?? "Failed to load venues.");
  }

  const body = (await response.json()) as ApiResponse<Venue[]>;
  return body.data;
}

function toSearchParams(params: VenueSearchParams) {
  const query = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "" || value === false) {
      return;
    }
    query.set(key, String(value));
  });

  const value = query.toString();
  return value ? `?${value}` : "";
}

function filterFallbackVenues(params: VenueSearchParams) {
  const q = params.q?.trim().toLowerCase();
  const tags = params.tags?.split(",").map((tag) => tag.trim()).filter(Boolean) ?? [];

  return fallbackVenues.filter((venue) => {
    if (q) {
      const searchable = [
        venue.name,
        venue.short_description,
        venue.about,
        venue.district,
        ...venue.categories,
        ...venue.trending_dishes
      ].join(" ").toLowerCase();
      if (!searchable.includes(q)) {
        return false;
      }
    }
    if (params.district && venue.district !== params.district) {
      return false;
    }
    if (params.max_price_vnd && venue.avg_price_max_vnd > params.max_price_vnd) {
      return false;
    }
    if (tags.length > 0 && !tags.some((tag) => venue.categories.includes(tag))) {
      return false;
    }
    return true;
  });
}
