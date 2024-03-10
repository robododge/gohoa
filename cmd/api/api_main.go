package main

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/robododge/gohoa"
)

func main() {

	gohoa.LoadStreetMappingsJson()
	r := gin.Default()
	r.GET("/ping/:aname", pongRespond)
	r.GET("/ping", pongRespond)

	apiController := NewApiController()
	apiController.registerRoutes(r)

	r.SetTrustedProxies(nil)

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

type ApiController struct {
	ml *gohoa.MemberLookup
}

func NewApiController() *ApiController {
	memberLookup := gohoa.NewMemberLookup()
	return &ApiController{memberLookup}
}

func (ac *ApiController) registerRoutes(r *gin.Engine) {
	r.GET("/suggest/streetnumber/:number", ac.suggestStreetNumber)
	r.GET("/suggest/streetname/:name", ac.suggestStreetName)
	r.GET("/suggest/streets_by_num/:number", ac.findMembersByStreetNumber)
}

func (ac *ApiController) suggestStreetNumber(c *gin.Context) {
	number := c.Param("number")
	ginH := gin.H{"response-type": "suggested street number"}
	if number != "" {
		ginH["data"] = ac.ml.SuggestNumber(number)
	}
	c.JSON(http.StatusOK, ginH)
}

func (ac *ApiController) suggestStreetName(c *gin.Context) {
	name := c.Param("name")
	ginH := gin.H{"response-type": "suggested street name"}
	if name != "" {
		ginH["data"] = ac.ml.SuggestStreetName(name)
	}
	c.JSON(http.StatusOK, ginH)
}

func (ac *ApiController) findMembersByStreetNumber(c *gin.Context) {
	number := c.Param("number")
	ginH := gin.H{"response-type": "members by street number"}
	if number != "" {
		if data, err := ac.ml.FindMembersByStreetNumber(number); err == nil {
			ginH["data"] = data
		} else {
			ginH["error"] = err
			c.JSON(http.StatusInternalServerError, ginH)
			return
		}
	}
	c.JSON(http.StatusOK, ginH)
}
