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

func GetUserByCollageName() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var collage struct {
			CollageName string `json:"collage_name" validate:"required,min=2,max=100"`
		}
		if err := c.BindJSON(&collage); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		validationErr := validate.Struct(collage)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		var users []models.LinkedInProfile

		filter := bson.D{{"education.collagename", collage.CollageName}}

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
		defer cursor.Close(ctx)
		if len(users) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No users found with the provided Collage"})
			return
		}

		c.JSON(http.StatusOK, users)

	}
}
