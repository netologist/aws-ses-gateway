package internal

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"net/http"
	"net/url"
	"strconv"
)

type RequestBody struct {
	Action string `json:"Action"`
}

func handler(c *gin.Context) {
	var reqBody RequestBody

	// Read the request body as a string
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	bodyString := string(bodyBytes)

	values, _ := url.ParseQuery(bodyString)
	reqBody = RequestBody{Action: values.Get("Action")}

	log.Println(reqBody) // prints the decoded request

	// Actions
	switch reqBody.Action {
	case "SendEmail":
		mailErr := SendEmail(bodyString)

		if mailErr != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": mailErr.Error(),
			})
			return
		}

		successTemplate, err := os.ReadFile("../assets/templates/success.xml")
		if err != nil {
			logrus.Error("Cannot open template success file: ", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
	
		// Replace {{message}} with absolute path of the body.html
		successMessage := strings.Replace(string(successTemplate), "{{message}}", "Mail sent", -1)
	
		// Respond with the content & 200
		c.String(http.StatusOK, successMessage)	
		break

	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unsupported action"})
		return
	}
}

func StartServer() {
	// Read environment variables
	ReadConfigFromEnv()
	logrus.Info("Starting mock server under port ", Config.Port)

	// Endpoints
	r := gin.Default()
	r.POST("/", handler)

	// Run
	err := r.Run(":" + strconv.Itoa(Config.Port))
	if err != nil {
		panic(err)
	}
}
