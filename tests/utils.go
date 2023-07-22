package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/wenzzyx/go-ushort/app/modules/auth"
)

func authFn(r *gin.Engine, authRoute, email, password string) (*auth.AuthResponse, error) {
	w := httptest.NewRecorder()

	reqBody, _ := json.Marshal(map[string]any{
		"email":    email,
		"password": password,
	})

	req, _ := http.NewRequest("POST", "/api/v1/auth"+authRoute, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resBody *auth.AuthResponse

	err := json.Unmarshal(w.Body.Bytes(), &resBody)
	return resBody, err
}

// Run this on top function of test file
func RegisterForTest(r *gin.Engine, email, password string) (*auth.AuthResponse, error) {
	return authFn(r, "/registration", email, password)
}

func LoginForTest(r *gin.Engine, email, password string) (*auth.AuthResponse, error) {
	return authFn(r, "/login", email, password)
}
