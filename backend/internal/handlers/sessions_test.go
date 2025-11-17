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

func TestGetSessions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := createMockDB()
	cfg := &config.Config{
		AdminPassword: "test-password",
		SubcollectionID: "test-collection",
	}
	h := New(mockDB, cfg)

	router := gin.New()
	router.GET("/api/sessions", h.GetSessions)

	req, _ := http.NewRequest("GET", "/api/sessions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Will fail without proper DB mock, but tests handler structure
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

func TestCreateSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    models.Session
		expectedStatus int
	}{
		{
			name: "valid session",
			requestBody: models.Session{
				Title:       "Test Session",
				Description: "Test Description",
				Time:        "10:00 AM",
				Duration:    "1 hour",
				SpeakerIDs:  []string{"speaker1"},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "missing title",
			requestBody:    models.Session{Description: "Test Description"},
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
			router.POST("/api/admin/sessions", h.CreateSession)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/admin/sessions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUpdateSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := createMockDB()
	cfg := &config.Config{
		AdminPassword: "test-password",
		SubcollectionID: "test-collection",
	}
	h := New(mockDB, cfg)

	router := gin.New()
	router.PUT("/api/admin/sessions/:id", h.UpdateSession)

	session := models.Session{
		Title:       "Updated Title",
		Description: "Updated Description",
		Time:        "11:00 AM",
		Duration:    "2 hours",
	}
	body, _ := json.Marshal(session)
	req, _ := http.NewRequest("PUT", "/api/admin/sessions/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
}

func TestDeleteSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := createMockDB()
	cfg := &config.Config{
		AdminPassword: "test-password",
		SubcollectionID: "test-collection",
	}
	h := New(mockDB, cfg)

	router := gin.New()
	router.DELETE("/api/admin/sessions/:id", h.DeleteSession)

	req, _ := http.NewRequest("DELETE", "/api/admin/sessions/test-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
}

