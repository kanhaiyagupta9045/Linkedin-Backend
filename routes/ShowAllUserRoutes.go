package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kanhaiyaKumarGupta/social-media/controllers/User"
)

func GetAllUser(router *gin.Engine) {
	router.GET("/linkedin-profiles/byName", User.GetAllUserByName())
	router.GET("/linkedin-profiles/bylocation", User.GetAllUserByLocation())
	router.GET("/linkedin-profiles/bySkills", User.GetUserBySkill())
	router.GET("/linkedin-profiles/byPosition", User.GetUserByPosition())
	router.GET("/linkedin-profiles/byCompany", User.GetUserByCompany())
	router.GET("/linkedin-profiles/byCollageName", User.GetUserByCollageName())
}
