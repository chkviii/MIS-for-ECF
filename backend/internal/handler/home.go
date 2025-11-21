package handler

import (
	"mypage-backend/internal/util"
	
	"github.com/gin-gonic/gin"
)

func HomeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Serve the index.html file
		c.File(util.Html_Path("index.html"))
	}
}