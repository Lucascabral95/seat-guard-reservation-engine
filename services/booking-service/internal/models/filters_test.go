package models_test

import (
	"testing"

	"booking-service/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedEvents(t *testing.T, db *gorm.DB) {
	events := []models.Event{
		{
			BaseModel: models.BaseModel{ID: uuid.NewString()},
			Name:      "Lollapalooza",
			Gender:    string(models.Rock),
			Location:  "Buenos Aires",
			Price:     50000,
		},
		{
			BaseModel: models.BaseModel{ID: uuid.NewString()},
			Name:      "Ultra Music Festival",
			Gender:    string(models.Electronica),
			Location:  "Miami",
			Price:     80000,
		},
		{
			BaseModel: models.BaseModel{ID: uuid.NewString()},
			Name:      "Buenos Aires Jazz",
			Gender:    string(models.Jazz),
			Location:  "Buenos Aires",
			Price:     20000,
		},
	}
	for _, e := range events {
		require.NoError(t, db.Create(&e).Error)
	}
}

func applyEventFilter(db *gorm.DB, f models.EventFilter) *gorm.DB {
	q := db.Model(&models.Event{})
	if f.Name != "" {
		q = q.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.Gender != "" {
		q = q.Where("gender = ?", f.Gender)
	}
	if f.Location != "" {
		q = q.Where("location = ?", f.Location)
	}
	return q
}

func TestEventFilter_ByName(t *testing.T) {
    db := setupTestDB(t)
    seedEvents(t, db)

    filter := models.EventFilter{Name: "Jazz"}

    var results []models.Event
    err := applyEventFilter(db, filter).Find(&results).Error

    assert.NoError(t, err)
    assert.Len(t, results, 1)
    assert.Equal(t, "Buenos Aires Jazz", results[0].Name)
}

func TestEventFilter_ByGender(t *testing.T) {
    db := setupTestDB(t)
    seedEvents(t, db)

    filter := models.EventFilter{Gender: string(models.Rock)}

    var results []models.Event
    err := applyEventFilter(db, filter).Find(&results).Error

    assert.NoError(t, err)
    assert.Len(t, results, 1)
    assert.Equal(t, "Lollapalooza", results[0].Name)
}

func TestEventFilter_ByLocation(t *testing.T) {
    db := setupTestDB(t)
    seedEvents(t, db)

    filter := models.EventFilter{Location: "Buenos Aires"}

    var results []models.Event
    err := applyEventFilter(db, filter).Find(&results).Error

    assert.NoError(t, err)
    assert.Len(t, results, 2)
}

func TestEventFilter_MultipleFields(t *testing.T) {
    db := setupTestDB(t)
    seedEvents(t, db)

    filter := models.EventFilter{
        Location: "Buenos Aires",
        Gender:   string(models.Rock),
    }

    var results []models.Event
    err := applyEventFilter(db, filter).Find(&results).Error

    assert.NoError(t, err)
    assert.Len(t, results, 1)
    assert.Equal(t, "Lollapalooza", results[0].Name)
}

func TestEventFilter_NoMatch(t *testing.T) {
    db := setupTestDB(t)
    seedEvents(t, db)

    filter := models.EventFilter{Name: "Tomorrowland"}

    var results []models.Event
    err := applyEventFilter(db, filter).Find(&results).Error

    assert.NoError(t, err)
    assert.Empty(t, results)
}

func TestEventFilter_Empty_ReturnsAll(t *testing.T) {
    db := setupTestDB(t)
    seedEvents(t, db)

    filter := models.EventFilter{}

    var results []models.Event
    err := applyEventFilter(db, filter).Find(&results).Error

    assert.NoError(t, err)
    assert.Len(t, results, 3)
}
