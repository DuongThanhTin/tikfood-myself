insert into venues (
  id,
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
  '11111111-1111-1111-1111-111111111111',
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
  '22222222-2222-2222-2222-222222222222',
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
on conflict (id) do nothing;

insert into dishes (
  id,
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
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
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
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2',
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
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1',
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
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2',
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
on conflict (id) do nothing;

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
) values
('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 35000, 60000, 'seed', 0.95, 24, 18, 120000, 92, now()),
('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', 30000, 55000, 'seed', 0.92, 18, 12, 80000, 86, now()),
('22222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 55000, 90000, 'seed', 0.94, 21, 15, 95000, 87, now()),
('22222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', 55000, 95000, 'seed', 0.88, 14, 9, 60000, 80, now())
on conflict (venue_id, dish_id) do nothing;

insert into tags (id, slug, label, type) values
('90000000-0000-0000-0000-000000000001', 'vietnamese', 'Vietnamese', 'cuisine'),
('90000000-0000-0000-0000-000000000002', 'street-food', 'Street food', 'venue_type'),
('90000000-0000-0000-0000-000000000003', 'breakfast', 'Breakfast', 'occasion'),
('90000000-0000-0000-0000-000000000004', 'late-night', 'Late night', 'occasion'),
('90000000-0000-0000-0000-000000000005', 'banh-mi', 'Banh mi', 'dish_type'),
('90000000-0000-0000-0000-000000000006', 'pho', 'Pho', 'dish_type')
on conflict (id) do nothing;

insert into venue_tags (venue_id, tag_id, confidence_score, source) values
('11111111-1111-1111-1111-111111111111', '90000000-0000-0000-0000-000000000001', 1, 'seed'),
('11111111-1111-1111-1111-111111111111', '90000000-0000-0000-0000-000000000002', 1, 'seed'),
('11111111-1111-1111-1111-111111111111', '90000000-0000-0000-0000-000000000004', 1, 'seed'),
('22222222-2222-2222-2222-222222222222', '90000000-0000-0000-0000-000000000001', 1, 'seed'),
('22222222-2222-2222-2222-222222222222', '90000000-0000-0000-0000-000000000003', 1, 'seed')
on conflict (venue_id, tag_id) do nothing;

insert into dish_tags (dish_id, tag_id, confidence_score, source) values
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', '90000000-0000-0000-0000-000000000005', 1, 'seed'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', '90000000-0000-0000-0000-000000000005', 1, 'seed'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', '90000000-0000-0000-0000-000000000006', 1, 'seed'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', '90000000-0000-0000-0000-000000000006', 1, 'seed')
on conflict (dish_id, tag_id) do nothing;

insert into ai_summaries (
  target_type,
  target_id,
  summary,
  model,
  prompt_version
) values
(
  'venue',
  '11111111-1111-1111-1111-111111111111',
  'Trending for late-night banh mi clips with strong local social proof.',
  'seed',
  'seed-v1'
),
(
  'venue',
  '22222222-2222-2222-2222-222222222222',
  'Popular for breakfast pho videos and consistent creator mentions.',
  'seed',
  'seed-v1'
);
