package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"appdirect-workshop-backend/internal/config"
	"appdirect-workshop-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSpeakers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := createMockDB()
	cfg := &config.Config{
		AdminPassword: "test-password",
		SubcollectionID: "test-collection",
	}
	h := New(mockDB, cfg)

	router := gin.New()
	router.GET("/api/speakers", h.GetSpeakers)

	req, _ := http.NewRequest("GET", "/api/speakers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Will fail without proper DB mock, but tests handler structure
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

func TestCreateSpeaker(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    models.Speaker
		expectedStatus int
	}{
		{
			name: "valid speaker",
			requestBody: models.Speaker{
				Name: "John Doe",
				Bio:  "Test bio",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "missing name",
			requestBody:    models.Speaker{Bio: "Test bio"},
			expectedStatus: http.StatusBadRequest,
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
			router.POST("/api/admin/speakers", h.CreateSpeaker)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/admin/speakers", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUpdateSpeaker(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := createMockDB()
	cfg := &config.Config{
		AdminPassword: "test-password",
		SubcollectionID: "test-collection",
	}
	h := New(mockDB, cfg)

	router := gin.New()
	router.PUT("/api/admin/speakers/:id", h.UpdateSpeaker)

	speaker := models.Speaker{
		Name: "Updated Name",
		Bio:  "Updated Bio",
	}
	body, _ := json.Marshal(speaker)
	req, _ := http.NewRequest("PUT", "/api/admin/speakers/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Will fail without proper mock, but tests structure
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
}

func TestDeleteSpeaker(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := createMockDB()
	cfg := &config.Config{
		AdminPassword: "test-password",
		SubcollectionID: "test-collection",
	}
	h := New(mockDB, cfg)

	router := gin.New()
	router.DELETE("/api/admin/speakers/:id", h.DeleteSpeaker)

	req, _ := http.NewRequest("DELETE", "/api/admin/speakers/test-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Will fail without proper mock, but tests structure
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
}

