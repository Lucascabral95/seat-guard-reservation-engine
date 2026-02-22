package services

import "testing"

func TestUpdateEventAvailability_EmptyEventIDReturnsNil(t *testing.T) {
	if err := UpdateEventAvailability(nil, ""); err != nil {
		t.Fatalf("expected nil error for empty eventID, got %v", err)
	}
}
