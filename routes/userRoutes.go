package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kanhaiyaKumarGupta/social-media/controllers"
)

func UserRoutes(router *gin.Engine) {

	router.POST("users/signup", controllers.SignUp())
	router.POST("users/login", controllers.Login())
	router.GET("users/getallusers", controllers.GetAllUsers())
	router.DELETE("user/:uid", controllers.DeleteUser())
	router.PUT("user/:uid/:password", controllers.UpdatePassword())
	router.POST("user/:uid", controllers.FindUser())
}
