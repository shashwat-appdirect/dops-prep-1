package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"appdirect-workshop-backend/internal/config"
	"appdirect-workshop-backend/internal/database"
	"appdirect-workshop-backend/internal/handlers"
	"appdirect-workshop-backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupTestRouter creates a test router with provided database interface
func SetupTestRouter(mockDB database.DBInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		SubcollectionID: "test-workshop",
		AdminPassword:    "test-password",
		Port:             "8080",
		CORSOrigin:      "http://localhost:5173",
		FirebaseServiceAccount: map[string]interface{}{
			"type": "service_account",
		},
	}

	h := handlers.New(mockDB, cfg)

	r := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.CORSOrigin}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/register", h.Register)
		public.GET("/registrations/count", h.GetRegistrationCount)
		public.GET("/speakers", h.GetSpeakers)
		public.GET("/sessions", h.GetSessions)
	}

	// Admin routes
	admin := r.Group("/api/admin")
	admin.POST("/login", h.AdminLogin)
	admin.Use(middleware.AuthMiddleware(cfg.AdminPassword))
	{
		admin.GET("/attendees", h.GetAttendees)
		admin.GET("/attendees/:id", h.GetAttendee)
		admin.GET("/speakers", h.GetSpeakers)
		admin.POST("/speakers", h.CreateSpeaker)
		admin.PUT("/speakers/:id", h.UpdateSpeaker)
		admin.DELETE("/speakers/:id", h.DeleteSpeaker)
		admin.GET("/sessions", h.GetSessions)
		admin.POST("/sessions", h.CreateSession)
		admin.PUT("/sessions/:id", h.UpdateSession)
		admin.DELETE("/sessions/:id", h.DeleteSession)
		admin.GET("/analytics/designations", h.GetDesignationBreakdown)
	}

	return r
}

// CreateTestRequest creates an HTTP test request with JSON body
func CreateTestRequest(method, url string, body interface{}) (*httptest.ResponseRecorder, *http.Request) {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return httptest.NewRecorder(), req
}

// GetAuthToken gets a test auth token from the router
func GetAuthToken(router *gin.Engine) (string, error) {
	w, req := CreateTestRequest("POST", "/api/admin/login", map[string]string{
		"password": "test-password",
	})
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return "", nil
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	return response["token"], nil
}
