package auth

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
	"github.com/wenzzyx/go-ushort/app/routers"
	"github.com/wenzzyx/go-ushort/tests"
)

var r *gin.Engine

const (
	basePath     = "/api/v1/auth"
	mockEmail    = "tester@tt.ttest"
	mockPassword = "qwerty1234"
)

func TestRegistration(t *testing.T) {
	cases := []struct {
		Email         string
		Password      string
		StatusCode    int
		MustSetCookie bool
	}{
		{
			Email:         mockEmail,
			Password:      mockPassword,
			StatusCode:    201,
			MustSetCookie: true,
		},
		{
			Email:         mockEmail,
			Password:      "",
			StatusCode:    422,
			MustSetCookie: false,
		},
		{
			Email:         mockEmail,
			Password:      "safsdfefsfs",
			StatusCode:    400,
			MustSetCookie: false,
		},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()

		reqBody, _ := json.Marshal(map[string]any{
			"email":    c.Email,
			"password": c.Password,
		})
		req, _ := http.NewRequest("POST", basePath+"/registration", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, c.StatusCode, w.Code)

		var resBody *any

		if c.MustSetCookie {
			err := json.Unmarshal(w.Body.Bytes(), &resBody)
			assert.Equal(t, nil, err)
			assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
			assert.Contains(t, *resBody, "accessToken")
		}

	}
}

func TestLogin(t *testing.T) {
	cases := []struct {
		Email         string
		Password      string
		StatusCode    int
		MustSetCookie bool
	}{
		{
			Email:         mockEmail,
			Password:      mockPassword,
			StatusCode:    201,
			MustSetCookie: true,
		},
		{
			Email:         mockEmail,
			Password:      "",
			StatusCode:    422,
			MustSetCookie: false,
		},
		{
			Email:         mockEmail,
			Password:      "safsdfefsfs",
			StatusCode:    400,
			MustSetCookie: false,
		},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()

		reqBody, _ := json.Marshal(map[string]any{
			"email":    c.Email,
			"password": c.Password,
		})
		req, _ := http.NewRequest("POST", basePath+"/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, c.StatusCode, w.Code)

		var resBody *any

		if c.MustSetCookie {
			err := json.Unmarshal(w.Body.Bytes(), &resBody)
			assert.Equal(t, nil, err)
			assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
			assert.Contains(t, *resBody, "accessToken")
		}

	}

}

func TestRefresh(t *testing.T) {
	w := httptest.NewRecorder()
	authParams, err := tests.LoginForTest(r, mockEmail, mockPassword)
	assert.Equal(t, nil, err)

	req, _ := http.NewRequest("GET", basePath+"/refresh", nil)
	req.Header.Set("Cookie", fmt.Sprintf("refreshToken=%s; Path=/", authParams.RefreshToken))
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var resBody *any

	err = json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, nil, err)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
	assert.Contains(t, *resBody, "accessToken")

}

func TestInvalidRefresh(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", basePath+"/refresh", nil)
	req.Header.Set("Cookie", fmt.Sprintf("refreshToken=%s; Path=/", "not-valit-refresh-token"))
	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

}

func TestLogout(t *testing.T) {
	w := httptest.NewRecorder()

	authParams, err := tests.LoginForTest(r, mockEmail, mockPassword)
	assert.Equal(t, nil, err)

	req, _ := http.NewRequest("DELETE", basePath+"/logout", nil)
	req.Header.Set("Cookie", fmt.Sprintf("refreshToken=%s; Path=/", authParams.RefreshToken))
	r.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}

func init() {
	tests.Setup()
	r = routers.SetupRouter()
}
