package handlers

import (
	"net/http"

	"appdirect-workshop-backend/internal/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handlers) GetSpeakers(c *gin.Context) {
	speakers := make([]models.Speaker, 0)
	iter := h.db.Collection("speakers").Documents(h.db.Context())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch speakers"})
			return
		}

		var speaker models.Speaker
		if err := doc.DataTo(&speaker); err != nil {
			continue
		}
		speaker.ID = doc.Ref.ID
		speakers = append(speakers, speaker)
	}

	// Always return an array, even if empty
	if speakers == nil {
		speakers = []models.Speaker{}
	}
	c.JSON(http.StatusOK, speakers)
}

func (h *Handlers) CreateSpeaker(c *gin.Context) {
	var speaker models.Speaker
	if err := c.ShouldBindJSON(&speaker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docRef, _, err := h.db.Collection("speakers").Add(h.db.Context(), speaker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create speaker"})
		return
	}

	speaker.ID = docRef.ID
	c.JSON(http.StatusCreated, speaker)
}

func (h *Handlers) UpdateSpeaker(c *gin.Context) {
	id := c.Param("id")
	var updates models.Speaker
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Collection("speakers").Doc(id).Set(h.db.Context(), updates)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Speaker not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update speaker"})
		return
	}

	updates.ID = id
	c.JSON(http.StatusOK, updates)
}

func (h *Handlers) DeleteSpeaker(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Collection("speakers").Doc(id).Delete(h.db.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Speaker not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete speaker"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Speaker deleted successfully"})
}

