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
  social_videos: SocialVideo[];
  trend_score: number;
  trending_dishes: string[];
  ai_summary: string;
  distance_meters?: number;
};

export type SocialVideo = {
  id: string;
  platform: "tiktok" | "instagram" | "youtube" | "facebook" | "other";
  url: string;
  creator_handle: string;
  caption: string;
  thumbnail_url: string;
  view_count: number;
  like_count: number;
  published_at?: string;
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
  city?: string;
  district?: string;
  dish?: string;
  tags?: string;
  platform?: string;
  lat?: number;
  lng?: number;
  radius_m?: number;
  min_price_vnd?: number;
  max_price_vnd?: number;
  open_now?: boolean;
  sort?: "trending" | "videos" | "distance" | "price";
  limit?: number;
};

export const fallbackVenues: Venue[] = [
  {
    id: "11111111-1111-4111-8111-111111111111",
    name: "Banh Mi Hem",
    slug: "banh-mi-hem-nguyen-trai-district-1",
    short_description: "Late-night banh mi spot trending on social video.",
    about: "A compact street-food venue known for grilled pork banh mi, quick service, and strong late-night local buzz near Nguyen Trai.",
    address: "12 Nguyen Trai",
    city: "Thành phố Hồ Chí Minh",
    district: "Quận 1",
    latitude: 10.7712,
    longitude: 106.6899,
    categories: ["banh-mi", "late-night", "street-food", "vietnamese"],
    price_level: 1,
    avg_price_min_vnd: 30000,
    avg_price_max_vnd: 80000,
    currency: "VND",
    social_video_count: 42,
    social_videos: [
      {
        id: "33333333-3333-4333-8333-333333333331",
        platform: "tiktok",
        url: "https://www.tiktok.com/@tikfood/video/banh-mi-hem-1",
        creator_handle: "@tikfood",
        caption: "Late-night banh mi with grilled pork near Nguyen Trai.",
        thumbnail_url: "",
        view_count: 120000,
        like_count: 8200,
        published_at: "2026-05-20T10:00:00Z"
      },
      {
        id: "33333333-3333-4333-8333-333333333332",
        platform: "instagram",
        url: "https://www.instagram.com/reel/banh-mi-hem-2/",
        creator_handle: "@saigonbites",
        caption: "Crispy banh mi pate and grilled pork combo.",
        thumbnail_url: "",
        view_count: 80000,
        like_count: 5100,
        published_at: "2026-05-22T10:00:00Z"
      }
    ],
    trend_score: 92,
    trending_dishes: ["banh mi thit nuong", "banh mi pate"],
    ai_summary: "Trending for late-night banh mi clips with strong local social proof."
  },
  {
    id: "22222222-2222-4222-8222-222222222222",
    name: "Pho Bo Nguyen",
    slug: "pho-bo-nguyen-le-van-sy-district-3",
    short_description: "Breakfast pho shop with consistent creator mentions.",
    about: "A neighborhood pho venue known for clear broth, beef toppings, and steady breakfast traffic from local regulars and food creators.",
    address: "88 Le Van Sy",
    city: "Thành phố Hồ Chí Minh",
    district: "Quận 3",
    latitude: 10.7864,
    longitude: 106.6767,
    categories: ["breakfast", "pho", "vietnamese"],
    price_level: 1,
    avg_price_min_vnd: 50000,
    avg_price_max_vnd: 120000,
    currency: "VND",
    social_video_count: 35,
    social_videos: [
      {
        id: "44444444-4444-4444-8444-444444444441",
        platform: "tiktok",
        url: "https://www.tiktok.com/@tikfood/video/pho-bo-nguyen-1",
        creator_handle: "@tikfood",
        caption: "Breakfast pho with clear broth and rare beef.",
        thumbnail_url: "",
        view_count: 95000,
        like_count: 6900,
        published_at: "2026-05-19T23:00:00Z"
      }
    ],
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
  const platforms = params.platform?.split(",").map((platform) => platform.trim()).filter(Boolean) ?? [];

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
    if (params.city && !matchesLocationAlias(venue.city, params.city)) {
      return false;
    }
    if (params.district && !matchesLocationAlias(venue.district, params.district)) {
      return false;
    }
    if (params.min_price_vnd && venue.avg_price_max_vnd < params.min_price_vnd) {
      return false;
    }
    if (params.max_price_vnd && venue.avg_price_max_vnd > params.max_price_vnd) {
      return false;
    }
    if (tags.length > 0 && !tags.some((tag) => venue.categories.includes(tag))) {
      return false;
    }
    if (platforms.length > 0 && !venue.social_videos.some((video) => platforms.includes(video.platform))) {
      return false;
    }
    return true;
  });
}

function matchesLocationAlias(value: string, target: string) {
  return normalizeLocationAlias(value) === normalizeLocationAlias(target);
}

function normalizeLocationAlias(value: string) {
  return value
    .trim()
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/\./g, "")
    .replace(/\s+/g, "")
    .replace("thanhphohochiminh", "hochiminh")
    .replace("tphochiminh", "hochiminh")
    .replace("tphcm", "hochiminh")
    .replace("hcm", "hochiminh")
    .replace("district1", "quan1")
    .replace("q1", "quan1")
    .replace("district3", "quan3")
    .replace("q3", "quan3");
}
