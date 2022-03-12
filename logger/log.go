package log

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// middelware
func DummyMiddleware(c *gin.Context) {
	c.Next()

	log.Println("Im a dummy!")

	message := fmt.Sprintf("\nWebsite:\n %v\n", (c.Request.Host + c.Request.URL.String()))
	//  msg += f'\nWebsite:\n{request._current_scheme_host}\n'

	message += fmt.Sprintf("\nError : \n %v : %v \n ", c.Writer.Status(), c.Errors.String())
	// msg += f'\nError:\n{err_type}: {self.exception}\n'

	message += fmt.Sprintf("\nMethod and Path: %v : %v\n", c.Request.Method, c.Request.URL.String())
	// msg += f'\nMethod & Path:\n[{request.method}] - {request.path}\n'

	message += fmt.Sprintf("\nPayload: %v \n", c.Request.Body)
	// msg += f'\nPayload:\n{payload}\n'

	message += fmt.Sprintf("\nParams: %v \n", c.Params)
	// msg += f'\nParams:\n{params}\n'

	// msg += f'\nUser:\n{user_info}\n'
	message += fmt.Sprintf("\nUser Info: %v \n", "user")

	// msg += f'\nStack:\n{self.stack}'
	message += fmt.Sprintf("\nstack: %v \n", c.Errors.Errors())

	log.Println(message)

	// Pass on to the next-in-chain
	statusCode := c.Writer.Status()
	log.Println(statusCode)
	if statusCode >= 500 {
		log.Println("error boos", statusCode)
	}
}
