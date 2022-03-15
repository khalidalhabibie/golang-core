package errors

import (
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/khalidalhabibie/golang-core/lib/dingtalk"
	"github.com/khalidalhabibie/golang-core/lib/redis"
)

// middelware
func ExceptionLoggingMiddleware(c *gin.Context) {
	c.Next()

	statusCode := c.Writer.Status()
	if statusCode >= 500 && os.Getenv("IS_REPORT_ERROR_ACTIVE") == "1" {
		// check, is exist in redis

		redisCredential := redis.Credentials{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		}

		redisClient := redis.NewClient(redisCredential, os.Getenv("APP_ENV"))

		// '{ACCESS_KEY}:{status_code}:{err_type}

		prefix := fmt.Sprintf("%s:", "exception_logging_middleware")
		key := fmt.Sprintf("%v:%v", statusCode, c.Errors.String())

		dataFromRedis := redisClient.Get(prefix, key)
		if dataFromRedis != "" {
			log.Println("error exist ", prefix+key)
			return
		}

		err := redisClient.Set(prefix, key, "exist", 1*time.Minute)
		if err != nil {
			log.Println("error set cache in the middleware: ", err)
			return
		}

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
