package errors

import (
	"fmt"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/khalidalhabibie/golang-core/lib/dingtalk"
)

// middelware
func PostErrorMiddleware(c *gin.Context) {
	c.Next()

	statusCode := c.Writer.Status()
	if statusCode >= 500 && os.Getenv("IS_REPORT_ERROR_ACTIVE") == "1" {
		claims := jwt.ExtractClaims(c)

		message := fmt.Sprintf("%v - ERROR %v ALERT\n", os.Getenv("APP_ENV"), statusCode)
		message += fmt.Sprintf("\nWebsite:\n%v\n", (c.Request.Host + c.Request.URL.String()))
		message += fmt.Sprintf("\nError : \n%v : %v \n ", statusCode, c.Errors.String())
		message += fmt.Sprintf("\nMethod and Path: %v : %v\n", c.Request.Method, c.Request.URL.String())
		message += fmt.Sprintf("\nPayload: \n%v \n", c.Request.PostForm)
		message += fmt.Sprintf("\nParams: %v \n", c.Request.URL.Query())
		message += fmt.Sprintf("\nUser Info:\n %v \n", claims)
		message += fmt.Sprintf("\nstack:\n %v \n", c.Errors.Errors())

		dingtalk.NewClient().SendtoDingtalk(os.Getenv("DINGTALK_ACCESS_KEY"),
			os.Getenv("DINGTALK_SECRET_KEY"),
			message, []string{}, []string{},
			true)

	}

}
