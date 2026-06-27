package openinghours

import (
	"testing"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

func byDay(hours []discovery.OpeningHour) map[int]discovery.OpeningHour {
	m := map[int]discovery.OpeningHour{}
	for _, h := range hours {
		m[h.DayOfWeek] = h
	}
	return m
}

func TestParseAllWeekSingleRange(t *testing.T) {
	hours := Parse("Mo-Su 07:00-23:00")
	if len(hours) != 7 {
		t.Fatalf("expected 7 days, got %d", len(hours))
	}
	m := byDay(hours)
	for day := 0; day <= 6; day++ {
		if m[day].OpenTime != "07:00" || m[day].CloseTime != "23:00" {
			t.Fatalf("day %d: got %+v", day, m[day])
		}
	}
}

func TestParseWeekdaysAndWeekend(t *testing.T) {
	hours := Parse("Mo-Fr 08:00-22:00; Sa-Su 09:00-23:00")
	m := byDay(hours)
	if m[5].OpenTime != "08:00" || m[5].CloseTime != "22:00" { // Friday
		t.Fatalf("friday: got %+v", m[5])
	}
	if m[0].OpenTime != "09:00" || m[0].CloseTime != "23:00" { // Sunday
		t.Fatalf("sunday: got %+v", m[0])
	}
	if m[6].OpenTime != "09:00" { // Saturday
		t.Fatalf("saturday: got %+v", m[6])
	}
}

func TestParseNoDayPrefixAppliesAllDays(t *testing.T) {
	hours := Parse("08:00-22:00")
	if len(hours) != 7 {
		t.Fatalf("expected 7 days, got %d", len(hours))
	}
}

func TestParseDayList(t *testing.T) {
	hours := Parse("Mo,We,Fr 08:00-18:00")
	m := byDay(hours)
	if len(hours) != 3 {
		t.Fatalf("expected 3 days, got %d", len(hours))
	}
	if _, ok := m[1]; !ok { // Monday
		t.Fatal("expected Monday")
	}
	if _, ok := m[2]; ok { // Tuesday should be absent
		t.Fatal("did not expect Tuesday")
	}
}

func TestParse247(t *testing.T) {
	hours := Parse("24/7")
	if len(hours) != 7 {
		t.Fatalf("expected 7 days, got %d", len(hours))
	}
	m := byDay(hours)
	if m[0].OpenTime != "00:00" || m[0].CloseTime != "23:59" {
		t.Fatalf("24/7 day: got %+v", m[0])
	}
}

func TestParseEmptyAndUnrecognized(t *testing.T) {
	if Parse("") != nil {
		t.Fatal("empty should be nil")
	}
	if Parse("sunrise-sunset") != nil {
		t.Fatal("unrecognized should be nil")
	}
}
