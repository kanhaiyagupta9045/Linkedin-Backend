package User

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kanhaiyaKumarGupta/social-media/databases"
	"github.com/kanhaiyaKumarGupta/social-media/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection = databases.GetCollection(databases.Client, "UserProfile")

var validate = validator.New()

func GetAllUserByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var searchInput struct {
			Name string `json:"name" validate:"required"`
		}
		if err := c.BindJSON(&searchInput); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		validationErr := validate.Struct(searchInput)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var users []models.LinkedInProfile
		filter := bson.D{{"firstname", searchInput.Name}}

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
			c.JSON(http.StatusNotFound, gin.H{"error": "No users found with the provided name"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}
