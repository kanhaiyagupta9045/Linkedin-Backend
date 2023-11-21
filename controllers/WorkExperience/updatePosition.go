package WorkExperience

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kanhaiyaKumarGupta/social-media/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdatePostion() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := c.Param("id")
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("Error converting ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		var updateposition struct {
			CompanyName string `json:"company_name" validate:"required,min=2,max=100"`
			Position    string `json:"position" validate:"required,min=2,max=100"`
		}
		if err := c.BindJSON(&updateposition); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		validationErr := validate.Struct(updateposition)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var user models.LinkedInProfile
		filter := bson.D{{"_id", objId}}

		err = collection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			log.Println("Error fetching existing user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing user data"})
			return
		}

		var companyFound bool
		for _, workExp := range user.WorkExperience {
			if workExp.CompanyName == updateposition.CompanyName {
				workExp.Position = updateposition.Position
				companyFound = true
				break
			}
		}
		if !companyFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found in user's work experience"})
			return
		}
		update := bson.D{
			{"$set", bson.D{
				{"workexperience", user.WorkExperience},
			}},
		}

		res, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println("Error updating Position:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Position"})
			return
		}

		if res.MatchedCount == 0 {
			log.Println("Document with ID not found:", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Position updated successfully"})
	}

}
