package controllers

import (
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
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var validate = validator.New()
var collection *mongo.Collection = databases.GetCollection(databases.Client, "UserProfile")

func UpdateUserSkills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("Error converting ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		var requestData struct {
			Skills []string `json:"skills"`
		}
		if err := c.BindJSON(&requestData); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var updateSkillsPointers []*string
		for _, skill := range requestData.Skills {
			skillCopy := skill
			updateSkillsPointers = append(updateSkillsPointers, &skillCopy)
		}

		filter := bson.D{{"_id", objID}}
		update := bson.D{
			{"$set", bson.D{
				{"skills", updateSkillsPointers},
			}},
		}
		upsert := false
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		res, err := collection.UpdateOne(ctx, filter, update, &opt)
		if err != nil {
			log.Println("Error updating skills:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update skills"})
			return
		}
		if res.MatchedCount == 0 {
			log.Println("Document with ID not found:", objID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Skills added successfully"})
	}
}

func AddSkills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		objID := c.Param("id")
		id, err := primitive.ObjectIDFromHex(objID)
		if err != nil {
			log.Println("Error converting ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		var requestskills struct {
			Skills []string `json:"skills"`
		}
		if err := c.BindJSON(&requestskills); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(requestskills)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		var addSkillsPointers []*string
		for _, skill := range requestskills.Skills {
			skillCopy := skill
			addSkillsPointers = append(addSkillsPointers, &skillCopy)
		}

		var user models.LinkedInProfile
		filter := bson.D{{"_id", id}}
		err = collection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			log.Println("Error fetching existing user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing user data"})
			return
		}
		for _, newSkill := range addSkillsPointers {
			exists := false

			for _, existingSkill := range user.Skills {
				existing := *existingSkill
				if *newSkill == existing {
					log.Println("Skill already exists:", *newSkill)
					exists = true
					break
				}
			}

			if !exists {
				skillCopy := *newSkill
				user.Skills = append(user.Skills, &skillCopy)
			}
		}

		update := bson.D{
			{"$set", bson.D{
				{"skills", user.Skills},
			}},
		}
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, err := collection.UpdateOne(ctx, filter, update, &opt)
		if err != nil {
			log.Println("Error updating skills:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update skills"})
			return
		}
		if res.MatchedCount == 0 {
			log.Println("Document with ID not found:", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Skills updated successfully"})
	}
}

func DeleteSkills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("Error converting ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		var skills struct {
			Skills []string `json:"skills"`
		}
		if err := c.BindJSON(&skills); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(skills)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		filter := bson.D{{"_id", objID}}

		// Fetch user data from the database
		var user models.LinkedInProfile
		err = collection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			log.Println("Error fetching existing user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing user data"})
			return
		}

		// Delete matching skills from the user's skills
		for _, deleteSkill := range skills.Skills {
			for i, existingSkill := range user.Skills {
				if *existingSkill == deleteSkill {
					log.Println("Skill to be deleted:", deleteSkill)
					// Remove the skill from user.Skills
					user.Skills = append(user.Skills[:i], user.Skills[i+1:]...)
					break
				}
			}
		}
		update := bson.D{
			{"$set", bson.D{
				{"skills", user.Skills},
			}},
		}

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		res, err := collection.UpdateOne(ctx, filter, update, &opt)
		if err != nil {
			log.Println("Error updating skills:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skills"})
			return
		}

		if res.MatchedCount == 0 {
			log.Println("Document with ID not found:", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Skills deleted successfully"})
	}
}

func ShowAllSkills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("Error converting ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		var user models.LinkedInProfile
		filter := bson.D{{"_id", objID}}
		err = collection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			log.Println("Error fetching existing user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing user data"})
			return
		}
		var skills []string
		for _, skill := range user.Skills {
			skills = append(skills, *skill)
		}
		c.JSON(http.StatusOK, skills)
	}
}
