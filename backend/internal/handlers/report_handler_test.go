package handlers

import (
	"database/sql"
	"testing"
	"time"
)

func TestCalculateOverviewPeriodStartUsesRequestedStartWhenNoServiceExists(t *testing.T) {
	requestedStart := time.Date(2026, time.August, 12, 0, 0, 0, 0, time.UTC)

	periodStart := calculateOverviewPeriodStart(requestedStart, sql.NullTime{})

	if !periodStart.Equal(requestedStart) {
		t.Fatalf("expected requested start %s, got %s", requestedStart, periodStart)
	}
}

func TestCalculateOverviewPeriodStartUsesNewestLowerBound(t *testing.T) {
	requestedStart := time.Date(2026, time.August, 12, 0, 0, 0, 0, time.UTC)
	serviceCreatedAt := time.Date(2026, time.August, 18, 21, 30, 0, 0, time.UTC)

	periodStart := calculateOverviewPeriodStart(requestedStart, sql.NullTime{
		Time:  serviceCreatedAt,
		Valid: true,
	})

	if !periodStart.Equal(serviceCreatedAt) {
		t.Fatalf("expected service creation time %s, got %s", serviceCreatedAt, periodStart)
	}
}

func TestCalculateOverviewPeriodStartKeepsFullWindowForOlderServices(t *testing.T) {
	requestedStart := time.Date(2026, time.August, 12, 0, 0, 0, 0, time.UTC)
	serviceCreatedAt := time.Date(2026, time.July, 18, 21, 30, 0, 0, time.UTC)

	periodStart := calculateOverviewPeriodStart(requestedStart, sql.NullTime{
		Time:  serviceCreatedAt,
		Valid: true,
	})

	if !periodStart.Equal(requestedStart) {
		t.Fatalf("expected requested start %s, got %s", requestedStart, periodStart)
	}
}
