package controllers

import (
	"context"
	"fmt"
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
)

var validate = validator.New()
var collection *mongo.Collection = databases.GetCollection(databases.Client, "UserProfile")

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err := c.BindJSON(&user)
		if err != nil {
			log.Fatal(err)
		}
		validationerr := validate.Struct(user)
		if validationerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erorr": validationerr.Error()})
			return
		}
		count, err := collection.CountDocuments(ctx, bson.M{"email": user.Email})
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email already registered"})
			return
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		result, err := collection.InsertOne(ctx, user)
		if err != nil {
			msg := fmt.Sprintf("user not created")
			log.Fatal(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result.InsertedID)
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var founduser models.User

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if founduser.Email == "" {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = collection.FindOne(ctx, bson.M{"user_id": founduser.User_id}).Decode(&founduser)

		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, founduser)

	}

}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		cursor, err := collection.Find(ctx, bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}

		var users []primitive.M
		for cursor.Next(ctx) {
			var user bson.M
			err := cursor.Decode(&user)
			if err != nil {
				log.Fatal(err)
			}

			users = append(users, user)
		}
		defer cancel()
		c.JSON(http.StatusOK, users)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user_id := c.Param("uid")

		id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		filter := bson.M{"_id": id}

		res, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userID := c.Param("uid")
		pwd := c.Param("password")

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			log.Printf("[%s] Error binding JSON: %v", time.Now().Format(time.RFC3339), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updatedObj primitive.D

		if user.Name != "" {
			updatedObj = append(updatedObj, bson.E{"name", user.Name})
		}
		if user.Email != "" {
			updatedObj = append(updatedObj, bson.E{"email", user.Email})
		}

		if pwd != "" {
			updatedObj = append(updatedObj, bson.E{"password", pwd})
		}

		user.Updated_at = time.Now()

		upsert := true

		filter := bson.M{"uid": userID}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		res, err := collection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updatedObj},
			},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprintf("[%s] User update failed: %v", time.Now().Format(time.RFC3339), err)
			log.Println(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "modifiedCount": res.ModifiedCount})
	}
}

func FindUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID := c.Param("uid")
		var user models.User

		id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}
		}

		filter := bson.M{"_id": id}
		err = collection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, user)
	}
}
