package handlers

import (
	"net/http"
	"time"

	"appdirect-workshop-backend/internal/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Designation string `json:"designation" binding:"required"`
}

func (h *Handlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create registration
	reg := models.Registration{
		Name:        req.Name,
		Email:       req.Email,
		Designation: req.Designation,
		CreatedAt:   time.Now(),
	}

	docRef, _, err := h.db.Collection("registrations").Add(h.db.Context(), reg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create registration"})
		return
	}

	reg.ID = docRef.ID
	c.JSON(http.StatusCreated, reg)
}

func (h *Handlers) GetRegistrationCount(c *gin.Context) {
	iter := h.db.Collection("registrations").Documents(h.db.Context())
	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count registrations"})
			return
		}
		count++
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

