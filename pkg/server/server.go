package server

import (
	"errors"
	"net/http"

	"github.com/Ski-Nav/SkiNav-Server/pkg/common/maps"
	"github.com/Ski-Nav/SkiNav-Server/pkg/lib/db"
	"github.com/gin-gonic/gin"
)

func Init() {
	db := db.Init()
	resortMap := maps.Init(db)
	r := gin.Default()
	r.GET("/maps", func(c *gin.Context) {
		c.JSON(http.StatusOK, resortMap.GetAllResorts())
	})
	r.GET("/maps/:ResortName", func(c *gin.Context) {
		resortName := c.Param("ResortName")
		graph, err := resortMap.GetGraphByResortName(resortName)
		if err != nil {
			c.AbortWithError(400, errors.New("resort not found"))
			return
		}
		c.JSON(http.StatusOK, graph)
	})
	r.Run(":3000")
}
