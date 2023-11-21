package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kanhaiyaKumarGupta/social-media/routes"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	routes.LinkedinRoutes(router)
	routes.Profile(router)
	routes.GetAllUser(router)
	routes.Experience(router)

	fmt.Println("Router is running on port:", port)

	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
