insert into venues (
  name,
  slug,
  short_description,
  about,
  address,
  city,
  district,
  latitude,
  longitude,
  location,
  price_level,
  avg_price_min_vnd,
  avg_price_max_vnd,
  social_video_count,
  social_post_count,
  total_view_count,
  total_like_count,
  source
) values
(
  'Banh Mi Hem',
  'banh-mi-hem-nguyen-trai-district-1',
  'Late-night banh mi spot trending on social video.',
  'A compact street-food venue known for grilled pork banh mi, quick service, and strong late-night local buzz near Nguyen Trai.',
  '12 Nguyen Trai',
  'Ho Chi Minh City',
  'District 1',
  10.7712,
  106.6899,
  ST_SetSRID(ST_MakePoint(106.6899, 10.7712), 4326)::geography,
  1,
  30000,
  80000,
  42,
  58,
  250000,
  18000,
  'seed'
),
(
  'Pho Bo Nguyen',
  'pho-bo-nguyen-le-van-sy-district-3',
  'Breakfast pho shop with consistent creator mentions.',
  'A neighborhood pho venue known for clear broth, beef toppings, and steady breakfast traffic from local regulars and food creators.',
  '88 Le Van Sy',
  'Ho Chi Minh City',
  'District 3',
  10.7864,
  106.6767,
  ST_SetSRID(ST_MakePoint(106.6767, 10.7864), 4326)::geography,
  1,
  50000,
  120000,
  35,
  44,
  180000,
  13500,
  'seed'
)
on conflict (slug) do update
set
  name = excluded.name,
  short_description = excluded.short_description,
  about = excluded.about,
  address = excluded.address,
  city = excluded.city,
  district = excluded.district,
  latitude = excluded.latitude,
  longitude = excluded.longitude,
  location = excluded.location,
  price_level = excluded.price_level,
  avg_price_min_vnd = excluded.avg_price_min_vnd,
  avg_price_max_vnd = excluded.avg_price_max_vnd,
  social_video_count = excluded.social_video_count,
  social_post_count = excluded.social_post_count,
  total_view_count = excluded.total_view_count,
  total_like_count = excluded.total_like_count,
  source = excluded.source,
  updated_at = now();

insert into dishes (
  name,
  normalized_name,
  slug,
  short_description,
  about,
  category,
  cuisine,
  typical_price_min_vnd,
  typical_price_max_vnd
) values
(
  'Banh mi thit nuong',
  'banh mi thit nuong',
  'banh-mi-thit-nuong',
  'Vietnamese grilled pork banh mi.',
  'A crispy baguette filled with grilled pork, pickles, herbs, pate, and sauces. It is a common social-food discovery item because each venue has a distinct style.',
  'sandwich',
  'vietnamese',
  30000,
  60000
),
(
  'Banh mi pate',
  'banh mi pate',
  'banh-mi-pate',
  'Classic pate-forward banh mi.',
  'A classic Vietnamese banh mi variation centered on rich pate, herbs, pickled vegetables, and a crisp baguette.',
  'sandwich',
  'vietnamese',
  25000,
  55000
),
(
  'Pho bo tai',
  'pho bo tai',
  'pho-bo-tai',
  'Beef pho with rare beef slices.',
  'A Vietnamese noodle soup with aromatic broth, rice noodles, herbs, and rare beef slices.',
  'noodle',
  'vietnamese',
  50000,
  90000
),
(
  'Pho bo vien',
  'pho bo vien',
  'pho-bo-vien',
  'Beef pho with beef balls.',
  'A beef pho variation known for springy beef balls, warm broth, and a filling breakfast profile.',
  'noodle',
  'vietnamese',
  50000,
  95000
)
on conflict (slug) do update
set
  name = excluded.name,
  normalized_name = excluded.normalized_name,
  short_description = excluded.short_description,
  about = excluded.about,
  category = excluded.category,
  cuisine = excluded.cuisine,
  typical_price_min_vnd = excluded.typical_price_min_vnd,
  typical_price_max_vnd = excluded.typical_price_max_vnd,
  updated_at = now();

insert into venue_dishes (
  venue_id,
  dish_id,
  price_min_vnd,
  price_max_vnd,
  price_source,
  confidence_score,
  mention_count,
  video_count,
  view_count,
  trend_score,
  last_seen_at
)
select
  venue.id,
  dish.id,
  seed.price_min_vnd,
  seed.price_max_vnd,
  'seed',
  seed.confidence_score,
  seed.mention_count,
  seed.video_count,
  seed.view_count,
  seed.trend_score,
  now()
from (
  values
    ('banh-mi-hem-nguyen-trai-district-1', 'banh-mi-thit-nuong', 35000, 60000, 0.95::numeric, 24, 18, 120000::bigint, 92::numeric),
    ('banh-mi-hem-nguyen-trai-district-1', 'banh-mi-pate', 30000, 55000, 0.92::numeric, 18, 12, 80000::bigint, 86::numeric),
    ('pho-bo-nguyen-le-van-sy-district-3', 'pho-bo-tai', 55000, 90000, 0.94::numeric, 21, 15, 95000::bigint, 87::numeric),
    ('pho-bo-nguyen-le-van-sy-district-3', 'pho-bo-vien', 55000, 95000, 0.88::numeric, 14, 9, 60000::bigint, 80::numeric)
) as seed(venue_slug, dish_slug, price_min_vnd, price_max_vnd, confidence_score, mention_count, video_count, view_count, trend_score)
join venues venue on venue.slug = seed.venue_slug
join dishes dish on dish.slug = seed.dish_slug
on conflict (venue_id, dish_id) do update
set
  price_min_vnd = excluded.price_min_vnd,
  price_max_vnd = excluded.price_max_vnd,
  price_source = excluded.price_source,
  confidence_score = excluded.confidence_score,
  mention_count = excluded.mention_count,
  video_count = excluded.video_count,
  view_count = excluded.view_count,
  trend_score = excluded.trend_score,
  last_seen_at = excluded.last_seen_at,
  updated_at = now();

insert into tags (slug, label, type) values
('vietnamese', 'Vietnamese', 'cuisine'),
('street-food', 'Street food', 'venue_type'),
('breakfast', 'Breakfast', 'occasion'),
('late-night', 'Late night', 'occasion'),
('banh-mi', 'Banh mi', 'dish_type'),
('pho', 'Pho', 'dish_type')
on conflict (slug) do update
set
  label = excluded.label,
  type = excluded.type;

insert into venue_tags (venue_id, tag_id, confidence_score, source)
select venue.id, tag.id, 1, 'seed'
from (
  values
    ('banh-mi-hem-nguyen-trai-district-1', 'vietnamese'),
    ('banh-mi-hem-nguyen-trai-district-1', 'street-food'),
    ('banh-mi-hem-nguyen-trai-district-1', 'late-night'),
    ('pho-bo-nguyen-le-van-sy-district-3', 'vietnamese'),
    ('pho-bo-nguyen-le-van-sy-district-3', 'breakfast')
) as seed(venue_slug, tag_slug)
join venues venue on venue.slug = seed.venue_slug
join tags tag on tag.slug = seed.tag_slug
on conflict (venue_id, tag_id) do update
set
  confidence_score = excluded.confidence_score,
  source = excluded.source;

insert into dish_tags (dish_id, tag_id, confidence_score, source)
select dish.id, tag.id, 1, 'seed'
from (
  values
    ('banh-mi-thit-nuong', 'banh-mi'),
    ('banh-mi-pate', 'banh-mi'),
    ('pho-bo-tai', 'pho'),
    ('pho-bo-vien', 'pho')
) as seed(dish_slug, tag_slug)
join dishes dish on dish.slug = seed.dish_slug
join tags tag on tag.slug = seed.tag_slug
on conflict (dish_id, tag_id) do update
set
  confidence_score = excluded.confidence_score,
  source = excluded.source;

insert into ai_summaries (
  target_type,
  target_id,
  summary,
  model,
  prompt_version
)
select
  'venue',
  venue.id,
  seed.summary,
  'seed',
  'seed-v1'
from (
  values
    ('banh-mi-hem-nguyen-trai-district-1', 'Trending for late-night banh mi clips with strong local social proof.'),
    ('pho-bo-nguyen-le-van-sy-district-3', 'Popular for breakfast pho videos and consistent creator mentions.')
) as seed(venue_slug, summary)
join venues venue on venue.slug = seed.venue_slug
where not exists (
  select 1
  from ai_summaries existing
  where existing.target_type = 'venue'
    and existing.target_id = venue.id
    and existing.prompt_version = 'seed-v1'
);

insert into venue_opening_hours (
  venue_id,
  day_of_week,
  open_time,
  close_time,
  is_closed
)
select
  venue.id,
  days.day_of_week,
  seed.open_time,
  seed.close_time,
  false
from (
  values
    ('banh-mi-hem-nguyen-trai-district-1', time '08:00', time '23:30'),
    ('pho-bo-nguyen-le-van-sy-district-3', time '06:00', time '14:00')
) as seed(venue_slug, open_time, close_time)
join venues venue on venue.slug = seed.venue_slug
cross join generate_series(0, 6) as days(day_of_week)
where not exists (
  select 1
  from venue_opening_hours existing
  where existing.venue_id = venue.id
    and existing.day_of_week = days.day_of_week
);
