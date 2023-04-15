package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Ski-Nav/SkiNav-Server/pkg/common/maps"
	"github.com/Ski-Nav/SkiNav-Server/pkg/lib/db"
	"github.com/gin-gonic/gin"
)

func Init() {
	db := db.Init()
	maps := maps.Init()
	r := gin.Default()
	r.GET("/maps/", func(c *gin.Context) {
		resortName := c.DefaultQuery("resort", "Big Bear")
		_, ok := maps[resortName]
		if !ok {
			c.AbortWithError(400, errors.New("no resort found"))
			return
		}
		fmt.Println(resortName)
		graph := db.GetGraphByResort(resortName)
		// jsonString, err := json.Marshal(graph)
		// if err != nil {
		// 	fmt.Println("Error marshaling graph to JSON:", err)
		// 	return
		// }
		c.JSON(http.StatusOK, graph)
	})
	r.Run(":3000")
}
