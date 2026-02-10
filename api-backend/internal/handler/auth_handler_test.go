package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginRequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "empty body",
			body:       map[string]string{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       map[string]string{"username": "admin"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing username",
			body:       map[string]string{"password": "123456"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			var req LoginRequest
			err := c.ShouldBindJSON(&req)
			if err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}

func TestLoginRequestValid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{
		"username": "admin",
		"password": "admin123",
	})
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Username != "admin" {
		t.Errorf("Username = %s, want admin", req.Username)
	}
	if req.Password != "admin123" {
		t.Errorf("Password = %s, want admin123", req.Password)
	}
}
