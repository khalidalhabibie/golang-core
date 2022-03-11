package log

import (
	"log"

	"github.com/gin-gonic/gin"
)

// middelware
func DummyMiddleware(c *gin.Context) {
	c.Next()

	log.Println("Im a dummy!")

	// Pass on to the next-in-chain
	statusCode := c.Writer.Status()
	log.Println(statusCode)
	if statusCode >= 500 {
		log.Println("error boos", statusCode)
	}
}
