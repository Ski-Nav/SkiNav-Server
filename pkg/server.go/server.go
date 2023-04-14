package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init() {
	r := gin.Default()
	r.GET("/maps", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"resort": "big bear",
		})
	})
	r.Run(":3000")
}
