package routes

import (
	"github.com/gin-gonic/gin"
	c "github.com/kanhaiyaKumarGupta/social-media/controllers/Profile"
)

func Profile(router *gin.Engine) {
	router.POST("/linkedin-profiles", c.CreateLinkedInProfile())
	router.POST("/linkedin-profiles/updateFirstName/:id", c.UpdateFirstName())
	router.POST("/linkedin-profiles/updateLastName/:id", c.UpdateLastName())
	router.POST("/linkedin-profiles/updatepassword/:id", c.UpdateUserPassword())
	router.POST("/linkedin-profiles/updatephone/:id", c.UpdatePhone())
	router.POST("/linkedin-profiles/updateEmail/:id", c.UpdateEmail())
	router.DELETE("/linkedin-profiles/deleteEmail/:id", c.DeleteLinkedInProfile())
}
