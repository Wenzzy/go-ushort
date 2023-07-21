package links

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-ushort/app/routers"
	"go-ushort/tests"
	"net/http"
	"net/http/httptest"
	_ "regexp"
	"testing"
)

var r *gin.Engine

var (
	basePath     = "/api/v1/links"
	mockEmail    = "tester@tt.ttest"
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

	authParams, err := tests.RegisterForTest(r, "test22231@tester.fsdf", mockPassword)
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
			fmt.Println(c, w.Code, w.Body.String())
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

func init() {
	tests.Setup()
	r = routers.SetupRouter()
}
