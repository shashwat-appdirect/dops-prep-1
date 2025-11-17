package config

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Save original env vars
	originalEnv := make(map[string]string)
	envVars := []string{
		"FIREBASE_SERVICE_ACCOUNT",
		"SUBSCOLLECTION_ID",
		"ADMIN_PASSWORD",
		"PORT",
		"CORS_ORIGIN",
	}
	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	tests := []struct {
		name          string
		setupEnv      func()
		expectedError bool
	}{
		{
			name: "valid config with base64 service account",
			setupEnv: func() {
				serviceAccount := map[string]interface{}{
					"type": "service_account",
					"project_id": "test-project",
				}
				jsonData, _ := json.Marshal(serviceAccount)
				base64Data := base64.StdEncoding.EncodeToString(jsonData)
				os.Setenv("FIREBASE_SERVICE_ACCOUNT", "base64:"+base64Data)
				os.Setenv("SUBSCOLLECTION_ID", "test-collection")
				os.Setenv("ADMIN_PASSWORD", "test-password")
			},
			expectedError: false,
		},
		{
			name: "missing FIREBASE_SERVICE_ACCOUNT (allowed for Cloud Run)",
			setupEnv: func() {
				os.Setenv("SUBSCOLLECTION_ID", "test-collection")
				os.Setenv("ADMIN_PASSWORD", "test-password")
				// Simulate Cloud Run environment
				os.Setenv("K_SERVICE", "test-service")
			},
			expectedError: false,
		},
		{
			name: "missing SUBSCOLLECTION_ID",
			setupEnv: func() {
				serviceAccount := map[string]interface{}{"type": "service_account"}
				jsonData, _ := json.Marshal(serviceAccount)
				base64Data := base64.StdEncoding.EncodeToString(jsonData)
				os.Setenv("FIREBASE_SERVICE_ACCOUNT", "base64:"+base64Data)
				os.Setenv("ADMIN_PASSWORD", "test-password")
			},
			expectedError: true,
		},
		{
			name: "missing ADMIN_PASSWORD",
			setupEnv: func() {
				serviceAccount := map[string]interface{}{"type": "service_account"}
				jsonData, _ := json.Marshal(serviceAccount)
				base64Data := base64.StdEncoding.EncodeToString(jsonData)
				os.Setenv("FIREBASE_SERVICE_ACCOUNT", "base64:"+base64Data)
				os.Setenv("SUBSCOLLECTION_ID", "test-collection")
			},
			expectedError: true,
		},
		{
			name: "default PORT",
			setupEnv: func() {
				serviceAccount := map[string]interface{}{"type": "service_account"}
				jsonData, _ := json.Marshal(serviceAccount)
				base64Data := base64.StdEncoding.EncodeToString(jsonData)
				os.Setenv("FIREBASE_SERVICE_ACCOUNT", "base64:"+base64Data)
				os.Setenv("SUBSCOLLECTION_ID", "test-collection")
				os.Setenv("ADMIN_PASSWORD", "test-password")
				os.Unsetenv("PORT")
			},
			expectedError: false,
		},
		{
			name: "default CORS_ORIGIN",
			setupEnv: func() {
				serviceAccount := map[string]interface{}{"type": "service_account"}
				jsonData, _ := json.Marshal(serviceAccount)
				base64Data := base64.StdEncoding.EncodeToString(jsonData)
				os.Setenv("FIREBASE_SERVICE_ACCOUNT", "base64:"+base64Data)
				os.Setenv("SUBSCOLLECTION_ID", "test-collection")
				os.Setenv("ADMIN_PASSWORD", "test-password")
				os.Unsetenv("CORS_ORIGIN")
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up env vars
			for _, key := range envVars {
				os.Unsetenv(key)
			}
			tt.setupEnv()

			cfg, err := Load()
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				if cfg != nil {
					assert.NotEmpty(t, cfg.SubcollectionID)
					assert.NotEmpty(t, cfg.AdminPassword)
					if os.Getenv("PORT") == "" {
						assert.Equal(t, "8080", cfg.Port)
					}
					if os.Getenv("CORS_ORIGIN") == "" {
						assert.Equal(t, "http://localhost:5173", cfg.CORSOrigin)
					}
				}
			}
		})
	}
}

