package storage

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestScheduledCampaign(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	var now = time.Now()

	store := From(db)

	segmentIDS := []int64{1, 2, 3, 4, 5, 6}
	segmentIDSsJSON, err := json.Marshal(segmentIDS)
	assert.Nil(t, err)

	cam := []*entities.Campaign{
		{
			UserID:       1,
			Name:         "test",
			TemplateID:   0,
			BaseTemplate: nil,
			Schedule:     nil,
			Status:       entities.StatusDraft,
		},
		{
			UserID:       1,
			Name:         "test2",
			TemplateID:   0,
			BaseTemplate: nil,
			Schedule:     nil,
			Status:       entities.StatusSending,
		},
		{
			UserID:       1,
			Name:         "test2",
			TemplateID:   0,
			BaseTemplate: nil,
			Schedule:     nil,
			Status:       entities.StatusScheduled,
		},
	}

	for i := range cam {
		err = store.CreateCampaign(cam[i])
		assert.Nil(t, err)
	}

	//Test create scheduled campaign
	c := []*entities.CampaignSchedule{
		{
			UserID:                  1,
			CampaignID:              cam[0].ID,
			ScheduledAt:             now,
			Source:                  "bla@email.com",
			FromName:                "from name",
			SegmentIDsJSON:          segmentIDSsJSON,
			DefaultTemplateDataJSON: []byte(`{"foo":"bar"}`),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			UserID:                  1,
			CampaignID:              cam[1].ID,
			ScheduledAt:             now,
			Source:                  "bla@email.com",
			FromName:                "from name",
			SegmentIDsJSON:          segmentIDSsJSON,
			DefaultTemplateDataJSON: []byte(`{"foo":"bar"}`),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			UserID:                  1,
			CampaignID:              cam[2].ID,
			ScheduledAt:             now,
			Source:                  "bla@email.com",
			FromName:                "from name",
			SegmentIDsJSON:          segmentIDSsJSON,
			DefaultTemplateDataJSON: []byte(`{"foo":"bar"}`),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
	}

	id := ksuid.New()
	for _, i := range c {
		i.ID = id
		err = store.CreateCampaignSchedule(i)
		assert.Nil(t, err)
		id = id.Next()
	}

	campSch, err := store.GetScheduledCampaigns(now)
	assert.Nil(t, err)

	// len should be 1 since the second campaign have status = sending (We only fetch campaigns with status draft and scheduled)
	assert.Equal(t, 2, len(campSch))

	// Test delete scheduled campaign
	err = store.DeleteCampaignSchedule(c[0].CampaignID)
	assert.Nil(t, err)
}
