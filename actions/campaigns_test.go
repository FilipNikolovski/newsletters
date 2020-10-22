package actions_test

import (
	"net/http"
	"testing"

	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/storage"
)

func TestCampaigns(t *testing.T) {
	s := storage.New("sqlite3", ":memory:")

	e := setup(t, s)
	auth, err := createAuthenticatedExpect(e, s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e.POST("/api/campaigns").WithForm(params.Campaign{Name: "djale", TemplateName: "djale"}).
		Expect().
		Status(http.StatusUnauthorized)

	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "", TemplateName: ""}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"name": "This field is required", "template_name": "This field is required"})

	// test post campaign
	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "djale", TemplateName: "djale"}).
		Expect().
		Status(http.StatusCreated)

	// test inserted campaign
	auth.GET("/api/campaigns/1").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "djale").
		ValueEqual("template_name", "djale").
		ValueEqual("status", "draft")

	auth.PUT("/api/campaigns/1").WithForm(params.Campaign{Name: "djaleputtest", TemplateName: "djaleputtest"}).
		Expect().
		Status(http.StatusNoContent)

	// test updated campaign
	auth.GET("/api/campaigns/1").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "djaleputtest").
		ValueEqual("template_name", "djaleputtest").
		ValueEqual("status", "draft")

	// delete campaign by id
	auth.DELETE("/api/campaigns/1").
		Expect().
		Status(http.StatusNoContent)
}