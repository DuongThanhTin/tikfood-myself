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
  ))
  and ($2 = '' or (
    v.district = $2
    or lower(unaccent(coalesce(v.district, ''))) = lower(unaccent($2))
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
group by v.id
order by
  case when $10 = 'distance' and $5::double precision is not null and $6::double precision is not null then ST_Distance(
    v.location,
    ST_SetSRID(ST_MakePoint($6::double precision, $5::double precision), 4326)::geography
  ) end asc nulls last,
  case when $10 = 'price' then v.avg_price_max_vnd end asc nulls last,
  case when $10 = 'videos' then v.social_video_count end desc nulls last,
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

func nullableFloat(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}
