package main

import (
	"log"
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/robododge/gohoa"
)

func main() {

	gohoa.LoadStreetMappingsJson()
	r := gin.Default()
	r.GET("/ping/:aname", pongRespond)
	r.GET("/ping", pongRespond)

	apiController := NewAPIController()
	apiController.registerRoutes(r)

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("something wrong with trusted proxi setting!", err)
	}

	err := r.Run() // listen and serve on
	if err != nil {
		log.Fatal(err)
	}
}

func pongRespond(c *gin.Context) {
	aname := c.Param("aname")

	ginH := gin.H{"message": "pong"}
	if aname != "" {
		ginH["name"] = aname
	}
	c.JSON(http.StatusOK, ginH)
}

type APIController struct {
	ml *gohoa.MemberLookup
	cr *gohoa.ContactReqService
}

func NewAPIController() *APIController {
	memberLookup := gohoa.NewMemberLookup()
	contactReqSvc := gohoa.NewContactReqService()
	return &APIController{memberLookup, contactReqSvc}
}

func (ac *APIController) registerRoutes(r *gin.Engine) {
	r.GET("/suggest/streetnumber/:number", ac.suggestStreetNumber)
	r.GET("/suggest/streetname/:name", ac.suggestStreetName)
	r.GET("/suggest/streets_by_num/:number", ac.findMembersByStreetNumber)
	r.POST("/contact/request", ac.createContactRequest)
}

func (ac *APIController) suggestStreetNumber(c *gin.Context) {
	number := c.Param("number")
	log.Println("go STD log.. logging from GIN: suggestStreetNumber: ", number)
	ginH := gin.H{"response-type": "suggested street number"}
	if number != "" {
		streetNumbers := ac.ml.SuggestNumber(number)
		log.Printf("Suggested street numbers for %s: len(%d)\n", number, len(streetNumbers))
		ginH["data"] = ac.ml.SuggestNumber(number)
	}
	c.JSON(http.StatusOK, ginH)
}

func (ac *APIController) suggestStreetName(c *gin.Context) {
	name := c.Param("name")
	ginH := gin.H{"response-type": "suggested street name"}
	if name != "" {
		ginH["data"] = ac.ml.SuggestStreetName(name)
	}
	c.JSON(http.StatusOK, ginH)
}

func (ac *APIController) findMembersByStreetNumber(c *gin.Context) {
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

func (ac *APIController) createContactRequest(c *gin.Context) {
	var contactReq gohoa.ContactRequest
	if err := c.BindJSON(&contactReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Contact request: %+v\n", contactReq)
	err := ac.cr.CreateContactRequest(contactReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Contact request received"})
}
