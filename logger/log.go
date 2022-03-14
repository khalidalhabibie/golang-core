package log

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	platform "github.com/khalidalhabibie/golang-core/platform"
)

// middelware
func DummyMiddleware(c *gin.Context) {
	c.Next()

	// Pass on to the next-in-chain
	statusCode := c.Writer.Status()
	if statusCode >= 500 && os.Getenv("IS_REPORT_ERROR_TO_DINGTALK") == "1" {
		message := fmt.Sprintf("%v - ERROR %v ALERT\n", os.Getenv("APP_ENV"), statusCode)

		message += fmt.Sprintf("\nWebsite:\n%v\n", (c.Request.Host + c.Request.URL.String()))
		//  msg += f'\nWebsite:\n{request._current_scheme_host}\n'

		message += fmt.Sprintf("\nError : \n%v : %v \n ", statusCode, c.Errors.String())
		// msg += f'\nError:\n{err_type}: {self.exception}\n'

		message += fmt.Sprintf("\nMethod and Path: %v : %v\n", c.Request.Method, c.Request.URL.String())
		// msg += f'\nMethod & Path:\n[{request.method}] - {request.path}\n'

		message += fmt.Sprintf("\nPayload: %v \n", c.Request)
		// msg += f'\nPayload:\n{payload}\n'

		message += fmt.Sprintf("\nParams: %v \n", c.Params)
		// msg += f'\nParams:\n{params}\n'

		// msg += f'\nUser:\n{user_info}\n'
		message += fmt.Sprintf("\nUser Info: %v \n", "user")

		// msg += f'\nStack:\n{self.stack}'
		message += fmt.Sprintf("\nstack: %v \n", c.Errors.Errors())

		platform.Dingtalk(os.Getenv("DINGTALK_ACCESS_KEY"),
			os.Getenv("DINGTALK_SECRET_KEY"),
			message, []string{}, []string{},
			true)

	}

}
