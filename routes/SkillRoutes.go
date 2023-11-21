package routes

import (
	"github.com/gin-gonic/gin"
	c "github.com/kanhaiyaKumarGupta/social-media/controllers/Skills"
)

func LinkedinRoutes(router *gin.Engine) {

	router.PUT("/linkedin-profiles/addskills/:id", c.AddSkills())
	router.DELETE("/linkedin-profiles/deleteskills/:id", c.DeleteSkills())
	router.GET("/linkedin-profiles/getAllSkills/:id", c.ShowAllSkills())
	router.PUT("/linkedin-profiles/:id", c.UpdateUserSkills())
	//router.GET("/linkedin-profiles", controllers.ListLinkedInProfiles())
	// router.GET("/linkedin-profiles/:id", controllers.GetLinkedInProfile())
	// router.DELETE("/linkedin-profiles/:id", controllers.DeleteLinkedInProfile())
	// router.PATCH("/linkedin-profiles/:id", controllers.PartialUpdateLinkedInProfile())
	// router.GET("/linkedin-profiles/by-email/:email", controllers.GetProfileByEmail())
	// router.GET("/linkedin-profiles/by-phone/:phone", controllers.GetProfileByPhone())
	// router.GET("/linkedin-profiles/by-location/:location", controllers.GetProfileByLocation())
	// router.GET("/linkedin-profiles/by-education/:education", controllers.GetProfileByEducation())
	// router.GET("/linkedin-profiles/by-work-experience/:work-experience", controllers.GetProfileByWorkExperience())
	// router.GET("/linkedin-profiles/by-skills/:skills", controllers.GetProfileBySkills())
	// router.GET("/linkedin-profiles/by-skill-and-location/:skill/:location", controllers.GetProfilesBySkillAndLocation())
	// router.GET("/linkedin-profiles/by-skill-and-education/:skill/:education", controllers.GetProfilesBySkillAndEducation())
	// router.GET("/linkedin-profiles/min-connections/:count", controllers.GetProfilesWithMinConnections())
	// router.GET("/linkedin-profiles/max-connections/:count", controllers.GetProfilesWithMaxConnections())
	// router.GET("/linkedin-profiles/min-max-connections/:min/:max", controllers.GetProfilesWithMinMaxConnections())
	// router.GET("linkedin-profiles/bycollageName/:collageName", controllers.GetProfileByCollageName())

}
