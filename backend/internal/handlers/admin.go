package handlers

import (
	"net/http"
	"time"

	"appdirect-workshop-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *Handlers) AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	if req.Password != h.cfg.AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.AdminPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

func (h *Handlers) GetAttendees(c *gin.Context) {
	var attendees []models.Registration
	iter := h.db.Collection("registrations").Documents(h.db.Context())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendees"})
			return
		}

		var reg models.Registration
		if err := doc.DataTo(&reg); err != nil {
			continue
		}
		reg.ID = doc.Ref.ID
		attendees = append(attendees, reg)
	}

	c.JSON(http.StatusOK, attendees)
}

func (h *Handlers) GetAttendee(c *gin.Context) {
	id := c.Param("id")
	doc, err := h.db.Collection("registrations").Doc(id).Get(h.db.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Attendee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendee"})
		return
	}

	var reg models.Registration
	if err := doc.DataTo(&reg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse attendee data"})
		return
	}
	reg.ID = doc.Ref.ID

	c.JSON(http.StatusOK, reg)
}

func (h *Handlers) GetDesignationBreakdown(c *gin.Context) {
	designationCount := make(map[string]int)
	iter := h.db.Collection("registrations").Documents(h.db.Context())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch registrations"})
			return
		}

		var reg models.Registration
		if err := doc.DataTo(&reg); err != nil {
			continue
		}
		designationCount[reg.Designation]++
	}

	var breakdown []models.DesignationBreakdown
	for designation, count := range designationCount {
		breakdown = append(breakdown, models.DesignationBreakdown{
			Designation: designation,
			Count:       count,
		})
	}

	c.JSON(http.StatusOK, breakdown)
}

