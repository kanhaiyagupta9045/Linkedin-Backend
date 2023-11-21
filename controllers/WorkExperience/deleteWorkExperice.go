package WorkExperience

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func DeleteWorkExperience() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	}
}
