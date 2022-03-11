package log

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// middelware
func DummyMiddleware(c *gin.Context) {
	c.Next()

	log.Println("Im a dummy!")

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Println("error load env ", err)
	// 	return
	// }

	log.Println("mesage core ", os.Getenv("MESSAGE"))

	// Pass on to the next-in-chain
	statusCode := c.Writer.Status()
	log.Println(statusCode)
	if statusCode >= 500 {
		log.Println("error boos", statusCode)
	}
}
