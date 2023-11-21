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

func CreateLinkedInProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var linkedinprofile models.LinkedInProfile
		if err := c.BindJSON(&linkedinprofile); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		if validationErr := validate.Struct(linkedinprofile); validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		filter := bson.D{{"email", linkedinprofile.Email}}
		count, err := collection.CountDocuments(ctx, filter)

		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email already registered"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email already registered"})
			return
		}

		result, err := collection.InsertOne(ctx, linkedinprofile)
		if err != nil {
			log.Println("Error inserting into the database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
			return
		}

		log.Println("Profile created successfully:", result.InsertedID)
		c.JSON(http.StatusOK, gin.H{"message": "Profile created successfully", "inserted_id": result.InsertedID})
	}
}

func UpdateUserPassword() gin.HandlerFunc {
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

		var newPassword struct {
			Password string `json:"password" validate:"required,min=8,max=50"`
		}
		if err := c.BindJSON(&newPassword); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		validationErr := validate.Struct(newPassword)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		filter := bson.D{{"_id", objID}}
		update := bson.D{
			{"$set", bson.D{
				{"password", newPassword.Password},
			}},
		}
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, err := collection.UpdateOne(ctx, filter, update, &opt)
		if err != nil {
			log.Println("Error updating password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		if res.MatchedCount == 0 {
			log.Println("Document with ID not found:", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}
}

func UpdatePhone() gin.HandlerFunc {
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

		var newPhone struct {
			Phone string `json:"phone" validate:"required,min=10,max=10"`
		}
		if err := c.BindJSON(&newPhone); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		validationErr := validate.Struct(newPhone)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		filter := bson.D{{"_id", objID}}
		update := bson.D{
			{"$set", bson.D{
				{"phone", newPhone.Phone}, // Update the "phone" field instead of "password"
			}},
		}

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		res, err := collection.UpdateOne(ctx, filter, update, &opt)
		if err != nil {
			log.Println("Error updating Phone:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update phone"})
			return
		}

		if res.MatchedCount == 0 {
			log.Println("Document with ID not found:", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Phone number updated successfully"})
	}
}

func UpdateEmail() gin.HandlerFunc {
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
		var newEmail struct {
			Email string `json:"email" validate:"email,required"`
		}
		if c.BindJSON(&newEmail); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		validationErr := validate.Struct(&newEmail)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		filter := bson.D{{"_id", objID}}

		update := bson.D{
			{"$set", bson.D{
				{"email", newEmail.Email},
			}},
		}
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		res, err := collection.UpdateOne(ctx, filter, update, &opt)
		if err != nil {
			log.Println("Error updating Email:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Email"})
			return
		}

		if res.MatchedCount == 0 {
			log.Println("Document with ID not found:", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully"})

	}
}

func DeleteLinkedInProfile() gin.HandlerFunc {
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

		filter := bson.D{{"_id", objID}}
		res, err := collection.DeleteOne(ctx, filter)

		if err != nil {
			log.Println("Error deleting document:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
			return
		}

		if res.DeletedCount == 0 {
			log.Println("Document with ID not found:", objID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})

	}
}

func UpdateFirstName() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		id := c.Param("id")
		ObjID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("Error converting ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		var ChangeFirstName struct {
			FirstName string `json:"first_name" validate:"required,min=6,max=50"`
		}

		if err := c.BindJSON(&ChangeFirstName); err != nil {

			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		validationErr := validate.Struct(&ChangeFirstName)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		filter := bson.D{{"_id", ObjID}}

		update := bson.D{{
			"$set", bson.D{
				{"firstname", ChangeFirstName.FirstName},
			}},
		}
		upsert := false
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, err := collection.UpdateOne(ctx, filter, update, &opt)
		if err != nil {
			log.Println("Error updating Email:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update First Name"})
			return
		}
		if res.MatchedCount == 0 {
			log.Println("Document with ID not found:", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "First Name updated successfully"})

	}
}

func UpdateLastName() gin.HandlerFunc {
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

		var changeLastName struct {
			LastName string `json:"last_name" validate:"required,min=6,max=15"`
		}
		if err := c.BindJSON(&changeLastName); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		validationErr := validate.Struct(changeLastName)
		if validationErr != nil {
			log.Println("Validation error:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		filter := bson.D{{"_id", objID}}
		update := bson.D{{
			"$set", bson.D{
				{"lastname", changeLastName.LastName},
			}},
		}
		upsert := false
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, err := collection.UpdateOne(ctx, filter, update, &opt)
		if err != nil {
			log.Println("Error updating LastName:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Last Name"})
			return
		}
		if res.MatchedCount == 0 {
			log.Println("Document with ID not found:", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with specified ID not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Last Name updated successfully"})
	}
}
