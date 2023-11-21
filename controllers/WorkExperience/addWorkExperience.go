package WorkExperience

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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection = databases.GetCollection(databases.Client, "UserProfile")
var validate = validator.New()

func AddWorkExperience() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := c.Param("id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("Error converting ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		var newWorkExperience struct {
			CompanyName string    `json:"company_name" validate:"required,min=2,max=100"`
			Position    string    `json:"position" validate:"required,min=2,max=100"`
			Description string    `json:"description" validate:"required,max=500"`
			StartDate   time.Time `json:"start_date"`
			EndDate     time.Time `json:"end_date"`
		}

		if err := c.BindJSON(&newWorkExperience); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		validationErr := validate.Struct(newWorkExperience)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var user models.LinkedInProfile
		filter := bson.D{{"_id", objID}}
		err = collection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			log.Fatal("No user found with corresponding id")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		for _, exp := range user.WorkExperience {
			if exp.CompanyName == newWorkExperience.CompanyName {
				log.Println("Company name already exists in user's work experience:", newWorkExperience.CompanyName)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Company name already exists in user's work experience"})
				return
			}
		}

		update := bson.D{
			{
				"$push", bson.D{
					{"workexperience", newWorkExperience},
				}},
		}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println("Error updating user's work experience:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update work experience"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Work experience added successfully"})

	}
}
