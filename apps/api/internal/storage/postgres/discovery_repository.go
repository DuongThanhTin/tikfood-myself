package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

type DiscoveryRepository struct {
	db *sql.DB
}

func NewDiscoveryRepository(db *sql.DB) *DiscoveryRepository {
	return &DiscoveryRepository{db: db}
}

func (repo *DiscoveryRepository) ListVenues(ctx context.Context, search discovery.VenueSearch) ([]discovery.Venue, error) {
	const query = `
select
  v.id::text,
  v.name,
  v.slug,
  coalesce(v.short_description, ''),
  coalesce(v.about, ''),
  coalesce(v.address_display, v.address, ''),
  coalesce((select city.official_name from locations city where city.id = v.city_location_id), v.city, ''),
  coalesce((select district.official_name from locations district where district.id = v.district_location_id), v.district, ''),
  coalesce(v.latitude, 0),
  coalesce(v.longitude, 0),
  coalesce((
    select json_agg(distinct tag_slug)
    from (
      select t.slug as tag_slug
      from venue_tags vt
      join tags t on t.id = vt.tag_id
      where vt.venue_id = v.id
      union
      select t.slug as tag_slug
      from venue_dishes vd_tag
      join dish_tags dt on dt.dish_id = vd_tag.dish_id
      join tags t on t.id = dt.tag_id
      where vd_tag.venue_id = v.id
    ) tag_list
  ), '[]'::json)::text as categories_json,
  coalesce(v.price_level, 0),
  coalesce(v.avg_price_min_vnd, 0),
  coalesce(v.avg_price_max_vnd, 0),
  coalesce(v.currency, 'VND'),
  coalesce(v.social_video_count, 0),
  coalesce((
    select json_agg(video_payload)
    from (
      select json_build_object(
        'id', sv.id::text,
        'platform', sv.platform,
        'url', sv.source_url,
        'creator_handle', coalesce(sv.creator_handle, ''),
        'caption', coalesce(sv.caption, ''),
        'thumbnail_url', coalesce(sv.thumbnail_url, ''),
        'view_count', coalesce(sv.view_count, 0),
        'like_count', coalesce(sv.like_count, 0),
        'published_at', case when sv.published_at is null then '' else to_char(sv.published_at at time zone 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') end
      ) as video_payload
      from social_videos sv
      where sv.venue_id = v.id
        and ($14 = '' or sv.platform = any(string_to_array($14, ',')))
      order by sv.view_count desc nulls last, sv.published_at desc nulls last
      limit 3
    ) ranked_videos
  ), '[]'::json)::text as social_videos_json,
  coalesce(max(vd.trend_score), 0)::int,
  coalesce(array_to_json(array_remove(array_agg(distinct d.name), null)), '[]'::json)::text as dishes_json,
  coalesce((
    select s.summary
    from ai_summaries s
    where s.target_type = 'venue'
      and s.target_id = v.id
    order by s.generated_at desc
    limit 1
  ), ''),
  case
    when $5::double precision is null or $6::double precision is null or v.location is null then null
    else ST_Distance(
      v.location,
      ST_SetSRID(ST_MakePoint($6::double precision, $5::double precision), 4326)::geography
    )
  end as distance_meters
from venues v
left join venue_dishes vd on vd.venue_id = v.id
left join dishes d on d.id = vd.dish_id
where ($1 = '' or (
    lower(v.name) like '%' || lower($1) || '%'
    or lower(coalesce(v.short_description, '')) like '%' || lower($1) || '%'
    or lower(coalesce(v.about, '')) like '%' || lower($1) || '%'
    or exists (
      select 1
      from venue_dishes vd_query
      join dishes d_query on d_query.id = vd_query.dish_id
      where vd_query.venue_id = v.id
        and (
          lower(d_query.name) like '%' || lower($1) || '%'
          or lower(d_query.normalized_name) like '%' || lower($1) || '%'
          or lower(coalesce(d_query.cuisine, '')) like '%' || lower($1) || '%'
          or lower(coalesce(d_query.category, '')) like '%' || lower($1) || '%'
        )
    )
    or exists (
      select 1
      from venue_tags vt_query
      join tags t_query on t_query.id = vt_query.tag_id
      where vt_query.venue_id = v.id
        and (lower(t_query.slug) like '%' || lower($1) || '%' or lower(t_query.label) like '%' || lower($1) || '%')
    )
    or exists (
      select 1
      from venue_dishes vd_tag_query
      join dish_tags dt_query on dt_query.dish_id = vd_tag_query.dish_id
      join tags t_query on t_query.id = dt_query.tag_id
      where vd_tag_query.venue_id = v.id
        and (lower(t_query.slug) like '%' || lower($1) || '%' or lower(t_query.label) like '%' || lower($1) || '%')
    )
    or exists (
      select 1
      from social_videos sv_query
      where sv_query.venue_id = v.id
        and (
          lower(coalesce(sv_query.caption, '')) like '%' || lower($1) || '%'
          or lower(coalesce(sv_query.creator_handle, '')) like '%' || lower($1) || '%'
          or lower(sv_query.platform) like '%' || lower($1) || '%'
        )
    )
  ))
  and ($12 = '' or (
    v.city = $12
    or lower(unaccent(coalesce(v.city, ''))) = lower(unaccent($12))
    or exists (
      select 1
      from locations city_location
      where city_location.id = v.city_location_id
        and (
          city_location.slug = lower(replace($12, ' ', '-'))
          or city_location.normalized_name = lower(unaccent($12))
        )
    )
    or exists (
      select 1
      from location_aliases city_alias
      where city_alias.location_id = v.city_location_id
        and city_alias.normalized_alias = lower(unaccent($12))
    )
  ))
  and ($2 = '' or (
    v.district = $2
    or lower(unaccent(coalesce(v.district, ''))) = lower(unaccent($2))
    or exists (
      select 1
      from locations district_location
      where district_location.id = v.district_location_id
        and (
          district_location.slug = lower(replace($2, ' ', '-'))
          or district_location.normalized_name = lower(unaccent($2))
        )
    )
    or exists (
      select 1
      from location_aliases district_alias
      where district_alias.location_id = v.district_location_id
        and district_alias.normalized_alias = lower(unaccent($2))
    )
  ))
  and ($3 = '' or exists (
    select 1
    from venue_dishes vd_filter
    join dishes d_filter on d_filter.id = vd_filter.dish_id
    where vd_filter.venue_id = v.id
      and (
        d_filter.normalized_name = lower($3)
        or lower(d_filter.name) like '%' || lower($3) || '%'
        or d_filter.slug = lower(replace($3, ' ', '-'))
      )
  ))
  and ($4 = '' or exists (
    select 1
    from (
      select t.slug
      from venue_tags vt
      join tags t on t.id = vt.tag_id
      where vt.venue_id = v.id
      union
      select t.slug
      from venue_dishes vd_tag
      join dish_tags dt on dt.dish_id = vd_tag.dish_id
      join tags t on t.id = dt.tag_id
      where vd_tag.venue_id = v.id
    ) filter_tags
    where filter_tags.slug = any(string_to_array($4, ','))
  ))
  and ($5::double precision is null or $6::double precision is null or $7 <= 0 or (
    v.location is not null
    and ST_DWithin(
      v.location,
      ST_SetSRID(ST_MakePoint($6::double precision, $5::double precision), 4326)::geography,
      $7
    )
  ))
  and ($8 <= 0 or (
    ($3 = '' and (v.avg_price_max_vnd is null or v.avg_price_max_vnd <= $8))
    or ($3 <> '' and exists (
      select 1
      from venue_dishes vd_price
      join dishes d_price on d_price.id = vd_price.dish_id
      where vd_price.venue_id = v.id
        and (
          d_price.normalized_name = lower($3)
          or lower(d_price.name) like '%' || lower($3) || '%'
          or d_price.slug = lower(replace($3, ' ', '-'))
        )
        and (vd_price.price_max_vnd is null or vd_price.price_max_vnd <= $8)
    ))
  ))
  and ($13 <= 0 or (
    ($3 = '' and (v.avg_price_max_vnd is null or v.avg_price_max_vnd >= $13))
    or ($3 <> '' and exists (
      select 1
      from venue_dishes vd_price_min
      join dishes d_price_min on d_price_min.id = vd_price_min.dish_id
      where vd_price_min.venue_id = v.id
        and (
          d_price_min.normalized_name = lower($3)
          or lower(d_price_min.name) like '%' || lower($3) || '%'
          or d_price_min.slug = lower(replace($3, ' ', '-'))
        )
        and (vd_price_min.price_max_vnd is null or vd_price_min.price_max_vnd >= $13)
    ))
  ))
  and (not $9::boolean or exists (
    select 1
    from venue_opening_hours oh
    where oh.venue_id = v.id
      and oh.day_of_week = extract(dow from now() at time zone 'Asia/Ho_Chi_Minh')::int
      and not oh.is_closed
      and oh.open_time is not null
      and oh.close_time is not null
      and (now() at time zone 'Asia/Ho_Chi_Minh')::time between oh.open_time and oh.close_time
  ))
  and ($14 = '' or exists (
    select 1
    from social_videos sv_platform
    where sv_platform.venue_id = v.id
      and sv_platform.platform = any(string_to_array($14, ','))
  ))
group by v.id
order by
  case when $10 = 'distance' and $5::double precision is not null and $6::double precision is not null then ST_Distance(
    v.location,
    ST_SetSRID(ST_MakePoint($6::double precision, $5::double precision), 4326)::geography
  ) end asc nulls last,
  case when $10 = 'price' then v.avg_price_max_vnd end asc nulls last,
  case when $10 = 'videos' then (
    select coalesce(sum(sv_rank.view_count), 0)
    from social_videos sv_rank
    where sv_rank.venue_id = v.id
      and ($14 = '' or sv_rank.platform = any(string_to_array($14, ',')))
  ) end desc nulls last,
  case when $10 = 'trending' then coalesce(max(vd.trend_score), 0) end desc nulls last,
  coalesce(max(vd.trend_score), 0) desc,
  v.social_video_count desc,
  v.name asc
limit $11
`

	rows, err := repo.db.QueryContext(
		ctx,
		query,
		search.Query,
		search.District,
		search.Dish,
		strings.Join(search.Tags, ","),
		nullableFloat(search.Lat),
		nullableFloat(search.Lng),
		search.RadiusM,
		search.MaxPriceVND,
		search.OpenNow,
		search.Sort,
		search.Limit,
		search.City,
		search.MinPriceVND,
		strings.Join(search.Platforms, ","),
	)
	if err != nil {
		return nil, fmt.Errorf("query venues: %w", err)
	}
	defer rows.Close()

	venues := []discovery.Venue{}
	for rows.Next() {
		var venue discovery.Venue
		var categoriesJSON string
		var dishesJSON string
		var socialVideosJSON string
		var distanceMeters sql.NullFloat64

		if err := rows.Scan(
			&venue.ID,
			&venue.Name,
			&venue.Slug,
			&venue.ShortDescription,
			&venue.About,
			&venue.Address,
			&venue.City,
			&venue.District,
			&venue.Latitude,
			&venue.Longitude,
			&categoriesJSON,
			&venue.PriceLevel,
			&venue.AvgPriceMinVND,
			&venue.AvgPriceMaxVND,
			&venue.Currency,
			&venue.SocialVideoCount,
			&socialVideosJSON,
			&venue.TrendScore,
			&dishesJSON,
			&venue.AISummary,
			&distanceMeters,
		); err != nil {
			return nil, fmt.Errorf("scan venue: %w", err)
		}

		if err := json.Unmarshal([]byte(categoriesJSON), &venue.Categories); err != nil {
			return nil, fmt.Errorf("decode venue categories: %w", err)
		}
		if err := json.Unmarshal([]byte(dishesJSON), &venue.TrendingDishes); err != nil {
			return nil, fmt.Errorf("decode venue dishes: %w", err)
		}
		if err := json.Unmarshal([]byte(socialVideosJSON), &venue.SocialVideos); err != nil {
			return nil, fmt.Errorf("decode venue social videos: %w", err)
		}
		if distanceMeters.Valid {
			venue.DistanceMeters = &distanceMeters.Float64
		}

		venues = append(venues, venue)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate venues: %w", err)
	}

	return venues, nil
}

func (repo *DiscoveryRepository) GetVenueBySlug(ctx context.Context, slug string) (discovery.Venue, error) {
	const query = `
select
  v.id::text,
  v.name,
  v.slug,
  coalesce(v.short_description, ''),
  coalesce(v.about, ''),
  coalesce(v.address_display, v.address, ''),
  coalesce((select city.official_name from locations city where city.id = v.city_location_id), v.city, ''),
  coalesce((select district.official_name from locations district where district.id = v.district_location_id), v.district, ''),
  coalesce(v.latitude, 0),
  coalesce(v.longitude, 0),
  coalesce((
    select json_agg(distinct tag_slug)
    from (
      select t.slug as tag_slug
      from venue_tags vt
      join tags t on t.id = vt.tag_id
      where vt.venue_id = v.id
      union
      select t.slug as tag_slug
      from venue_dishes vd_tag
      join dish_tags dt on dt.dish_id = vd_tag.dish_id
      join tags t on t.id = dt.tag_id
      where vd_tag.venue_id = v.id
    ) tag_list
  ), '[]'::json)::text as categories_json,
  coalesce(v.price_level, 0),
  coalesce(v.avg_price_min_vnd, 0),
  coalesce(v.avg_price_max_vnd, 0),
  coalesce(v.currency, 'VND'),
  coalesce(v.social_video_count, 0),
  coalesce((
    select json_agg(video_payload)
    from (
      select json_build_object(
        'id', sv.id::text,
        'platform', sv.platform,
        'url', sv.source_url,
        'creator_handle', coalesce(sv.creator_handle, ''),
        'caption', coalesce(sv.caption, ''),
        'thumbnail_url', coalesce(sv.thumbnail_url, ''),
        'view_count', coalesce(sv.view_count, 0),
        'like_count', coalesce(sv.like_count, 0),
        'published_at', case when sv.published_at is null then '' else to_char(sv.published_at at time zone 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') end
      ) as video_payload
      from social_videos sv
      where sv.venue_id = v.id
      order by sv.view_count desc nulls last, sv.published_at desc nulls last
      limit 10
    ) ranked_videos
  ), '[]'::json)::text as social_videos_json,
  coalesce((
    select max(vd_score.trend_score)
    from venue_dishes vd_score
    where vd_score.venue_id = v.id
  ), 0)::int,
  coalesce((
    select array_to_json(array_remove(array_agg(distinct d_names.name), null))
    from venue_dishes vd_names
    join dishes d_names on d_names.id = vd_names.dish_id
    where vd_names.venue_id = v.id
  ), '[]'::json)::text as trending_dishes_json,
  coalesce((
    select s.summary
    from ai_summaries s
    where s.target_type = 'venue'
      and s.target_id = v.id
    order by s.generated_at desc
    limit 1
  ), ''),
  coalesce((
    select json_agg(dish_payload order by (dish_payload->>'trend_score')::numeric desc nulls last, dish_payload->>'name')
    from (
      select json_build_object(
        'id', d.id::text,
        'name', d.name,
        'slug', d.slug,
        'short_description', coalesce(d.short_description, ''),
        'about', coalesce(d.about, ''),
        'category', coalesce(d.category, ''),
        'cuisine', coalesce(d.cuisine, ''),
        'price_min_vnd', coalesce(vd.price_min_vnd, d.typical_price_min_vnd, 0),
        'price_max_vnd', coalesce(vd.price_max_vnd, d.typical_price_max_vnd, 0),
        'currency', coalesce(vd.currency, d.currency, 'VND'),
        'mention_count', coalesce(vd.mention_count, 0),
        'video_count', coalesce(vd.video_count, 0),
        'view_count', coalesce(vd.view_count, 0),
        'trend_score', coalesce(vd.trend_score, 0)
      ) as dish_payload
      from venue_dishes vd
      join dishes d on d.id = vd.dish_id
      where vd.venue_id = v.id
    ) dish_rows
  ), '[]'::json)::text as dishes_json,
  coalesce((
    select json_agg(json_build_object(
      'day_of_week', oh.day_of_week,
      'open_time', case when oh.open_time is null then '' else to_char(oh.open_time, 'HH24:MI') end,
      'close_time', case when oh.close_time is null then '' else to_char(oh.close_time, 'HH24:MI') end,
      'is_closed', oh.is_closed
    ) order by oh.day_of_week)
    from venue_opening_hours oh
    where oh.venue_id = v.id
  ), '[]'::json)::text as opening_hours_json
from venues v
where v.slug = $1
limit 1
`

	var venue discovery.Venue
	var categoriesJSON string
	var socialVideosJSON string
	var trendingDishesJSON string
	var dishesJSON string
	var openingHoursJSON string

	err := repo.db.QueryRowContext(ctx, query, strings.TrimSpace(slug)).Scan(
		&venue.ID,
		&venue.Name,
		&venue.Slug,
		&venue.ShortDescription,
		&venue.About,
		&venue.Address,
		&venue.City,
		&venue.District,
		&venue.Latitude,
		&venue.Longitude,
		&categoriesJSON,
		&venue.PriceLevel,
		&venue.AvgPriceMinVND,
		&venue.AvgPriceMaxVND,
		&venue.Currency,
		&venue.SocialVideoCount,
		&socialVideosJSON,
		&venue.TrendScore,
		&trendingDishesJSON,
		&venue.AISummary,
		&dishesJSON,
		&openingHoursJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return discovery.Venue{}, discovery.ErrVenueNotFound
		}
		return discovery.Venue{}, fmt.Errorf("query venue detail: %w", err)
	}

	if err := json.Unmarshal([]byte(categoriesJSON), &venue.Categories); err != nil {
		return discovery.Venue{}, fmt.Errorf("decode venue categories: %w", err)
	}
	if err := json.Unmarshal([]byte(socialVideosJSON), &venue.SocialVideos); err != nil {
		return discovery.Venue{}, fmt.Errorf("decode venue social videos: %w", err)
	}
	if err := json.Unmarshal([]byte(trendingDishesJSON), &venue.TrendingDishes); err != nil {
		return discovery.Venue{}, fmt.Errorf("decode venue trending dishes: %w", err)
	}
	if err := json.Unmarshal([]byte(dishesJSON), &venue.Dishes); err != nil {
		return discovery.Venue{}, fmt.Errorf("decode venue dishes: %w", err)
	}
	if err := json.Unmarshal([]byte(openingHoursJSON), &venue.OpeningHours); err != nil {
		return discovery.Venue{}, fmt.Errorf("decode venue opening hours: %w", err)
	}

	return venue, nil
}

func nullableFloat(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}
