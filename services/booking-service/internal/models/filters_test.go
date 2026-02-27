package models_test

import (
	"testing"

	"booking-service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestEventFilter_ZeroValue(t *testing.T) {
	filter := models.EventFilter{}

	assert.Empty(t, filter.Name)
	assert.Empty(t, filter.Gender)
	assert.Empty(t, filter.Location)
}

func TestEventFilter_Assignment(t *testing.T) {
	filter := models.EventFilter{
		Name:     "Lollapalooza",
		Gender:   "ROCK",
		Location: "Buenos Aires",
	}

	assert.Equal(t, "Lollapalooza", filter.Name)
	assert.Equal(t, "ROCK", filter.Gender)
	assert.Equal(t, "Buenos Aires", filter.Location)
}
