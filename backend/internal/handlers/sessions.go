package handlers

import (
	"net/http"

	"appdirect-workshop-backend/internal/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handlers) GetSessions(c *gin.Context) {
	sessions := make([]models.Session, 0)
	iter := h.db.Collection("sessions").Documents(h.db.Context())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
			return
		}

		var session models.Session
		if err := doc.DataTo(&session); err != nil {
			continue
		}
		session.ID = doc.Ref.ID
		sessions = append(sessions, session)
	}

	// Always return an array, even if empty
	if sessions == nil {
		sessions = []models.Session{}
	}
	c.JSON(http.StatusOK, sessions)
}

func (h *Handlers) CreateSession(c *gin.Context) {
	var session models.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docRef, _, err := h.db.Collection("sessions").Add(h.db.Context(), session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	session.ID = docRef.ID
	c.JSON(http.StatusCreated, session)
}

func (h *Handlers) UpdateSession(c *gin.Context) {
	id := c.Param("id")
	var updates models.Session
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Collection("sessions").Doc(id).Set(h.db.Context(), updates)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	updates.ID = id
	c.JSON(http.StatusOK, updates)
}

func (h *Handlers) DeleteSession(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Collection("sessions").Doc(id).Delete(h.db.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session deleted successfully"})
}

