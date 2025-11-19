package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"appdirect-workshop-backend/internal/config"
	"appdirect-workshop-backend/internal/database"
	"appdirect-workshop-backend/internal/handlers"
	"appdirect-workshop-backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Firestore
	db, err := database.NewFirestoreClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer db.Close()

	// Set up Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	// Remove trailing slash from CORS origin if present
	corsOrigin := strings.TrimSuffix(cfg.CORSOrigin, "/")
	corsConfig.AllowOrigins = []string{corsOrigin}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// Initialize handlers
	h := handlers.New(db, cfg)

	// Serve static files (frontend build)
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		// Use StaticFS with http.Dir to serve static files with proper MIME types
		assetsPath := filepath.Join(staticDir, "assets")
		if _, err := os.Stat(assetsPath); err == nil {
			r.StaticFS("/assets", http.Dir(assetsPath))
			log.Printf("Serving assets from: %s", assetsPath)
			// Log assets directory contents for debugging
			if entries, err := os.ReadDir(assetsPath); err == nil {
				log.Printf("Found %d files in assets directory", len(entries))
			}
		} else {
			log.Printf("Warning: Assets directory not found at: %s", assetsPath)
		}
		r.StaticFile("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))
		r.StaticFile("/vite.svg", filepath.Join(staticDir, "vite.svg"))
		log.Printf("Static files directory found at: %s", staticDir)
	} else {
		log.Printf("Warning: Static files directory not found at: %s", staticDir)
	}

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

	// SPA routing fallback - serve index.html for non-API routes
	// This must be last to catch all non-API routes
	r.NoRoute(func(c *gin.Context) {
		// Don't serve index.html for API routes or asset routes
		if strings.HasPrefix(c.Request.URL.Path, "/api") || strings.HasPrefix(c.Request.URL.Path, "/assets") {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		// Serve index.html for all other routes (SPA routing)
		indexPath := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
		} else {
			log.Printf("Error: index.html not found at %s", indexPath)
			c.JSON(404, gin.H{"error": "Frontend not found"})
		}
	})

	// Start server
	// Cloud Run sets PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

