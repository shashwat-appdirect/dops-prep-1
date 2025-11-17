package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminPassword := "test-password-123"

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		shouldPass     bool
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + createTestToken(t, adminPassword),
			expectedStatus: http.StatusOK,
			shouldPass:     true,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			shouldPass:     false,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			shouldPass:     false,
		},
		{
			name:           "token without Bearer prefix",
			authHeader:     createTestToken(t, adminPassword),
			expectedStatus: http.StatusOK, // Middleware accepts token without Bearer prefix
			shouldPass:     true,
		},
		{
			name:           "expired token",
			authHeader:     "Bearer " + createExpiredToken(t, adminPassword),
			expectedStatus: http.StatusUnauthorized,
			shouldPass:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(AuthMiddleware(adminPassword))
			r.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.shouldPass {
				assert.Contains(t, w.Body.String(), "success")
			} else {
				assert.Contains(t, w.Body.String(), "error")
			}
		})
	}
}

func createTestToken(t *testing.T, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": true,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}
	return tokenString
}

func createExpiredToken(t *testing.T, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": true,
		"exp":   time.Now().Add(-24 * time.Hour).Unix(), // Expired
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}
	return tokenString
}

