package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/modak-test/api/controller"
	"github.com/jordanlanch/modak-test/bootstrap"
)

func NewNotificationRouter(env *bootstrap.Env, timeout time.Duration, rdb *redis.Client, group *gin.RouterGroup) {
	nc := controller.NewNotificationController(env, rdb)
	group.POST("/notification", nc.HandleNotification)
}
