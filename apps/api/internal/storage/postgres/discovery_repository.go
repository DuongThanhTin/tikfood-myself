package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

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
  coalesce(v.address, ''),
  coalesce(v.city, ''),
  coalesce(v.district, ''),
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
  ), '')
from venues v
left join venue_dishes vd on vd.venue_id = v.id
left join dishes d on d.id = vd.dish_id
where ($1 = '' or v.district = $1)
  and ($2 = '' or exists (
    select 1
    from venue_dishes vd_filter
    join dishes d_filter on d_filter.id = vd_filter.dish_id
    where vd_filter.venue_id = v.id
      and d_filter.normalized_name = lower($2)
  ))
group by v.id
order by v.social_video_count desc, v.name asc
`

	rows, err := repo.db.QueryContext(ctx, query, search.District, search.Dish)
	if err != nil {
		return nil, fmt.Errorf("query venues: %w", err)
	}
	defer rows.Close()

	venues := []discovery.Venue{}
	for rows.Next() {
		var venue discovery.Venue
		var categoriesJSON string
		var dishesJSON string

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
		); err != nil {
			return nil, fmt.Errorf("scan venue: %w", err)
		}

		if err := json.Unmarshal([]byte(categoriesJSON), &venue.Categories); err != nil {
			return nil, fmt.Errorf("decode venue categories: %w", err)
		}
		if err := json.Unmarshal([]byte(dishesJSON), &venue.TrendingDishes); err != nil {
			return nil, fmt.Errorf("decode venue dishes: %w", err)
		}

		venues = append(venues, venue)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate venues: %w", err)
	}

	return venues, nil
}
