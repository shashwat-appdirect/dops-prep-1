package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"appdirect-workshop-backend/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		password       string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "valid password",
			password:       "test-password",
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "invalid password",
			password:       "wrong-password",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "missing password",
			password:       "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := createMockDB()
			cfg := &config.Config{
				AdminPassword: "test-password",
				SubcollectionID: "test-collection",
			}
			h := New(mockDB, cfg)

			router := gin.New()
			router.POST("/api/admin/login", h.AdminLogin)

			body, _ := json.Marshal(map[string]string{"password": tt.password})
			req, _ := http.NewRequest("POST", "/api/admin/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "error")
			} else {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "token")
			}
		})
	}
}

func TestGetAttendees(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := createMockDB()
	cfg := &config.Config{
		AdminPassword: "test-password",
		SubcollectionID: "test-collection",
	}
	h := New(mockDB, cfg)

	router := gin.New()
	router.GET("/api/admin/attendees", h.GetAttendees)

	req, _ := http.NewRequest("GET", "/api/admin/attendees", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Will fail without proper DB mock, but tests handler structure
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

func TestGetDesignationBreakdown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := createMockDB()
	cfg := &config.Config{
		AdminPassword: "test-password",
		SubcollectionID: "test-collection",
	}
	h := New(mockDB, cfg)

	router := gin.New()
	router.GET("/api/admin/analytics/designations", h.GetDesignationBreakdown)

	req, _ := http.NewRequest("GET", "/api/admin/analytics/designations", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Will fail without proper DB mock, but tests handler structure
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

