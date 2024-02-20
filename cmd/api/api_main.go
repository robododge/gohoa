package main

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping/:aname", pongRespond)
	r.GET("/ping", pongRespond)
	r.Run() // listen and serve on
}

func pongRespond(c *gin.Context) {
	aname := c.Param("aname")

	ginH := gin.H{"message": "pong"}
	if aname != "" {
		ginH["name"] = aname
	}
	c.JSON(http.StatusOK, ginH)
}
