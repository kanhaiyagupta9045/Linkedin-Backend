package routes

import (
	"github.com/gin-gonic/gin"
	c "github.com/kanhaiyaKumarGupta/social-media/controllers/WorkExperience"
)

func Experience(router *gin.Engine) {
	router.POST("/linkedin-profiles/addwrokExperience/:id", c.AddWorkExperience())
	router.POST("/linkedin-profiles/updatepostion/:id", c.UpdatePostion())
}
