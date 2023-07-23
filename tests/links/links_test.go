package links

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	_ "regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	db_utils "github.com/wenzzyx/go-ushort/app/common/db-utils"
	"github.com/wenzzyx/go-ushort/app/routers"
	"github.com/wenzzyx/go-ushort/tests"
)

var r *gin.Engine

var (
	basePath     = "/api/v1/links"
	mockEmail    = "tester_links@tt.ttest"
	mockPassword = "qwerty1234"
)

var generatedAlias string

func TestCreateLink(t *testing.T) {
	cases := []struct {
		Name         string
		RealUrl      string
		StatusCode   int
		MustGetAlias bool
	}{
		{
			Name:         "Test link 1",
			RealUrl:      "https://google.com",
			StatusCode:   201,
			MustGetAlias: true,
		},
		{
			Name:         "Test link 1",
			RealUrl:      "",
			StatusCode:   422,
			MustGetAlias: false,
		},
		{
			Name:         "Test link 1",
			RealUrl:      "https://google.com",
			StatusCode:   200,
			MustGetAlias: true,
		},
	}

	authParams, err := tests.RegisterForTest(r, mockEmail, mockPassword)
	assert.Equal(t, nil, err)
	for _, c := range cases {
		w := httptest.NewRecorder()
		reqBody, _ := json.Marshal(map[string]any{
			"name":    c.Name,
			"realUrl": c.RealUrl,
		})

		req, _ := http.NewRequest("POST", basePath+"/", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authParams.AccessToken))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, c.StatusCode, w.Code)
		if c.MustGetAlias {
			var resBody *struct {
				ID    uint   `json:"id"`
				Alias string `json:"alias"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &resBody)
			assert.Equal(t, nil, err)
			assert.NotEmpty(t, resBody.Alias)
			generatedAlias = resBody.Alias

		}

	}
}

func TestUpdateLink(t *testing.T) {
	cases := []struct {
		ID                    int
		Name                  string
		StatusCode            int
		MustUseAuthentication bool
	}{
		{
			ID:                    1,
			Name:                  "Test-link-22",
			MustUseAuthentication: true,
			StatusCode:            204,
		},
		{
			ID:                    100,
			Name:                  "Test-link-22",
			MustUseAuthentication: true,
			StatusCode:            404,
		},
		{
			ID:                    100,
			Name:                  "Test-link-22",
			MustUseAuthentication: false,
			StatusCode:            401,
		},
	}

	authParams, err := tests.LoginForTest(r, mockEmail, mockPassword)
	assert.Equal(t, nil, err)
	for _, c := range cases {
		w := httptest.NewRecorder()
		reqBody, _ := json.Marshal(map[string]any{
			"name": c.Name,
		})

		req, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/%d", basePath, c.ID), bytes.NewBuffer(reqBody))
		if c.MustUseAuthentication {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authParams.AccessToken))
		}

		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, c.StatusCode, w.Code)
	}
}

func TestGetAllLinks(t *testing.T) {
	cases := []struct {
		StatusCode            int
		MustUseAuthentication bool
		MustCheckTotalCount   bool
	}{
		{
			StatusCode:            200,
			MustUseAuthentication: true,
			MustCheckTotalCount:   true,
		},
		{
			StatusCode:            401,
			MustUseAuthentication: false,
			MustCheckTotalCount:   false,
		},
	}

	authParams, err := tests.LoginForTest(r, mockEmail, mockPassword)
	assert.Equal(t, nil, err)
	for _, c := range cases {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", basePath+"/", nil)
		if c.MustUseAuthentication {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authParams.AccessToken))
		}
		r.ServeHTTP(w, req)
		assert.Equal(t, c.StatusCode, w.Code)
		if c.MustCheckTotalCount {
			var resBody *db_utils.GetAllResponse
			err := json.Unmarshal(w.Body.Bytes(), &resBody)
			assert.Equal(t, nil, err)
			fmt.Println(w.Body.String())
			assert.Greater(t, resBody.TotalCount, 0)

		}

	}
}

func TestRedirect(t *testing.T) {
	cases := []struct {
		Alias      string
		StatusCode int
	}{
		{
			Alias:      generatedAlias,
			StatusCode: 302,
		},
		{
			Alias:      generatedAlias + "22222",
			StatusCode: 404,
		},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/"+c.Alias, nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, c.StatusCode, w.Code)

	}
}

func init() {
	tests.Setup()
	r = routers.SetupRouter()
}
