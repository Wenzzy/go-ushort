package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-ushort/app/modules/auth"
	"net/http"
	"net/http/httptest"
)

func RegisterForTest(r *gin.Engine, email, password string) (*auth.AuthResponse, error) {
	w := httptest.NewRecorder()

	reqBody, _ := json.Marshal(map[string]any{
		"email":    email,
		"password": password,
	})

	req, _ := http.NewRequest("POST", "/api/v1/auth/registration", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resBody *auth.AuthResponse

	err := json.Unmarshal(w.Body.Bytes(), &resBody)
	return resBody, err
}
