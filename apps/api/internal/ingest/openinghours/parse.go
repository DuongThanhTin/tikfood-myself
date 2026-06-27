package openinghours

import (
	"regexp"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

// osmDays lists OSM weekday abbreviations in their natural order, paired with
// our day index (0 = Sunday ... 6 = Saturday).
var osmDays = []struct {
	abbr  string
	index int
}{
	{"Mo", 1}, {"Tu", 2}, {"We", 3}, {"Th", 4}, {"Fr", 5}, {"Sa", 6}, {"Su", 0},
}

var timeRangeRe = regexp.MustCompile(`(\d{2}:\d{2})\s*-\s*(\d{2}:\d{2})`)

func dayPosition(abbr string) int {
	for pos, d := range osmDays {
		if d.abbr == abbr {
			return pos
		}
	}
	return -1
}

// Parse converts a subset of OSM opening_hours syntax into per-day hours.
// Returns nil for empty or unrecognized input. Last rule wins per day.
func Parse(spec string) []discovery.OpeningHour {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil
	}
	if spec == "24/7" {
		return allDays("00:00", "23:59")
	}

	byDay := map[int]discovery.OpeningHour{}
	for _, rule := range strings.Split(spec, ";") {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		match := timeRangeRe.FindStringSubmatch(rule)
		if match == nil {
			continue
		}
		open, close := match[1], match[2]

		dayPart := strings.TrimSpace(rule[:strings.Index(rule, match[0])])
		days := parseDays(dayPart)
		for _, day := range days {
			byDay[day] = discovery.OpeningHour{
				DayOfWeek: day,
				OpenTime:  open,
				CloseTime: close,
				IsClosed:  false,
			}
		}
	}

	if len(byDay) == 0 {
		return nil
	}
	result := make([]discovery.OpeningHour, 0, len(byDay))
	for day := 0; day <= 6; day++ {
		if h, ok := byDay[day]; ok {
			result = append(result, h)
		}
	}
	return result
}

// parseDays expands a day specifier like "Mo-Fr", "Mo,We,Fr", or "" (all days).
func parseDays(spec string) []int {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		all := make([]int, 0, 7)
		for _, d := range osmDays {
			all = append(all, d.index)
		}
		return all
	}

	days := map[int]bool{}
	for _, token := range strings.Split(spec, ",") {
		token = strings.TrimSpace(token)
		if strings.Contains(token, "-") {
			parts := strings.SplitN(token, "-", 2)
			startPos := dayPosition(strings.TrimSpace(parts[0]))
			endPos := dayPosition(strings.TrimSpace(parts[1]))
			if startPos < 0 || endPos < 0 {
				continue
			}
			pos := startPos
			for {
				days[osmDays[pos].index] = true
				if pos == endPos {
					break
				}
				pos = (pos + 1) % len(osmDays)
			}
			continue
		}
		if p := dayPosition(token); p >= 0 {
			days[osmDays[p].index] = true
		}
	}

	out := make([]int, 0, len(days))
	for day := range days {
		out = append(out, day)
	}
	return out
}

func allDays(open string, close string) []discovery.OpeningHour {
	hours := make([]discovery.OpeningHour, 0, 7)
	for day := 0; day <= 6; day++ {
		hours = append(hours, discovery.OpeningHour{DayOfWeek: day, OpenTime: open, CloseTime: close})
	}
	return hours
}
