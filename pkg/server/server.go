package server

import (
	"errors"
	"net/http"

	"github.com/Ski-Nav/SkiNav-Server/pkg/common/maps"
	"github.com/Ski-Nav/SkiNav-Server/pkg/lib/db"
	"github.com/gin-gonic/gin"
)

func Init() {
	// initialize database connection
	db := db.Init()
	// init map
	resortMap := maps.Init(db)
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/app_compatibility", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		v1.GET("/maps", func(c *gin.Context) {
			c.JSON(http.StatusOK, resortMap.GetAllResorts())
		})
		v1.GET("/maps/:ResortName", func(c *gin.Context) {
			resortName := c.Param("ResortName")
			graph, err := resortMap.GetGraphByResortName(resortName)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, errors.New("resort not found"))
				return
			}
			c.JSON(http.StatusOK, graph)
		})
	}
	r.Run(":3000")
}
