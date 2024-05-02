package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/modak-test/bootstrap"
)

func Setup(env *bootstrap.Env, timeout time.Duration, rdb *redis.Client) *gin.Engine {
	router := gin.New()

	// All Public APIs
	publicRouter := router.Group("/api")
	NewNotificationRouter(env, timeout, rdb, publicRouter)

	return router
}
