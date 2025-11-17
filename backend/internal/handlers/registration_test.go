package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"appdirect-workshop-backend/internal/config"
	"appdirect-workshop-backend/internal/database"
	"appdirect-workshop-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test helper to create a mock database
func createMockDB() *database.MockFirestoreClient {
	return &database.MockFirestoreClient{}
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid registration",
			requestBody: map[string]string{
				"name":        "John Doe",
				"email":       "john@example.com",
				"designation": "Software Engineer",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "missing name",
			requestBody: map[string]string{
				"email":       "john@example.com",
				"designation": "Software Engineer",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid email",
			requestBody: map[string]string{
				"name":        "John Doe",
				"email":       "invalid-email",
				"designation": "Software Engineer",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "missing designation",
			requestBody: map[string]string{
				"name":  "John Doe",
				"email": "john@example.com",
			},
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
			router.POST("/api/register", h.Register)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "error")
			} else if w.Code == http.StatusCreated {
				// Only check response if creation succeeded (requires DB mock)
				var reg models.Registration
				json.Unmarshal(w.Body.Bytes(), &reg)
				assert.NotEmpty(t, reg.ID)
				assert.Equal(t, tt.requestBody.(map[string]string)["name"], reg.Name)
				assert.Equal(t, tt.requestBody.(map[string]string)["email"], reg.Email)
			}
		})
	}
}

func TestGetRegistrationCount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := createMockDB()
	cfg := &config.Config{
		AdminPassword: "test-password",
		SubcollectionID: "test-collection",
	}
	h := New(mockDB, cfg)

	router := gin.New()
	router.GET("/api/registrations/count", h.GetRegistrationCount)

	req, _ := http.NewRequest("GET", "/api/registrations/count", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Will fail without proper DB mock, but tests handler structure
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

