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

func GetUserByPosition() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var position struct {
			Position string `json:"position" validate:"required,min=2,max=100"`
		}
		if err := c.BindJSON(&position); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		validationErr := validate.Struct(position)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		filter := bson.D{{"workexperience.position", position.Position}}

		var users []models.LinkedInProfile
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			log.Println("Error fetching user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
			return
		}
		if err := cursor.All(ctx, &users); err != nil {
			log.Println("Error decoding user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user data"})
			return
		}
		if len(users) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No users found with the provided position"})
			return
		}
		c.JSON(http.StatusOK, users)

	}
}
