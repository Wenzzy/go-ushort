package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-ushort/app/modules/auth"
	"go-ushort/app/routers"
	"net/http"
	"net/http/httptest"
	_ "regexp"
	"testing"
)

func init() {
	Setup()
}

const basePath = "/api/v1/auth"

func TestRegistration(t *testing.T) {
	r := routers.SetupRouter()
	w := httptest.NewRecorder()

	reqBody, _ := json.Marshal(auth.RegisterValidator{
		Email:    "tester@tt.ttest",
		Password: "qwerty1234",
	})

	req, _ := http.NewRequest("POST", basePath+"/registration", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	fmt.Printf("Testssll: %v", w.Body.String())
	assert.Equal(t, 201, w.Code)

	var resBody *any

	err := json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, nil, err)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
	assert.Contains(t, *resBody, "accessToken")

}
