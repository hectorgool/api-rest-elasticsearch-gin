package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hectorgool/api-rest-elasticsearch-gin/common"
	"github.com/hectorgool/api-rest-elasticsearch-gin/elasticsearch"
)

func main() {

	r := gin.Default()
	r.Use(common.CORSMiddleware())

	// Ping test
	r.GET("/", func(c *gin.Context) {
		result, err := elasticsearch.Ping()
		common.CheckError(err)
		c.JSON(200, gin.H{"data": result})
	})

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	r.GET("/search", func(c *gin.Context) {
		q := c.DefaultQuery("q", "")
		result, err := elasticsearch.Search(q)
		common.CheckError(err)
		c.JSON(200, gin.H{"data": result})
	})

	r.Run() // listen and serve on 0.0.0.0:8080

}
