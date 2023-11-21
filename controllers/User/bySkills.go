package User

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kanhaiyaKumarGupta/social-media/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserBySkill() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var skillInput struct {
			Skill string `json:"skill" validate:"required"`
		}

		if err := c.BindJSON(&skillInput); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		// Assuming you have a collection named 'users' and 'skills' is an array field in the model
		filter := bson.D{{"skills", skillInput.Skill}}

		var users []models.LinkedInProfile
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			log.Println("Error fetching user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &users); err != nil {
			log.Println("Error decoding user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user data"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No users found with the provided skill"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}
