package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/invopop/jsonschema"
	"github.com/jordanlanch/modak-test/bootstrap"
	"github.com/jordanlanch/modak-test/domain"
	"github.com/jordanlanch/modak-test/repository"
	"github.com/jordanlanch/modak-test/usecase"
	"github.com/stoewer/go-strcase"
	"github.com/xeipuuv/gojsonschema"
)

type NotificationController struct {
	NotificationUseCase domain.NotificationService
}

// NewNotificationController inicializa un nuevo controlador de notificaciones con la lógica del caso de uso necesaria.
func NewNotificationController(env *bootstrap.Env, rdb *redis.Client) *NotificationController {
	repo := repository.NewNotificationRepository(rdb)
	limiter := usecase.NewRedisRateLimiter(rdb)
	notificationUseCase := usecase.NewNotificationUseCase(repo, limiter)

	return &NotificationController{
		NotificationUseCase: notificationUseCase,
	}
}

// HandleNotification procesa las solicitudes entrantes para enviar notificaciones, aplicando límites de tasa.
func (nc *NotificationController) HandleNotification(c *gin.Context) {
	payload := domain.Notification{}

	r := new(jsonschema.Reflector)
	r.KeyNamer = strcase.SnakeCase // from package github.com/stoewer/go-strcase
	// r.RequiredFromJSONSchemaTags = true
	payloadSchemaReflect := r.Reflect(domain.Notification{})
	schemaLoader := gojsonschema.NewGoLoader(&payloadSchemaReflect)

	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	documentLoader := gojsonschema.NewBytesLoader(requestBody)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if !result.Valid() {
		err = fmt.Errorf("[ERROR] invalid payload %+v", result.Errors())
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if err := json.Unmarshal(requestBody, &payload); err != nil {
		err = fmt.Errorf("[ERROR] error unmarshaling JSON %+v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})

		return
	}

	// Procesar la notificación
	ctx := c.Request.Context()
	err = nc.NotificationUseCase.SendNotification(ctx, payload)
	if err != nil {
		log.Printf("Error al enviar la notificación: %v", err)
		var httpStatus int
		if _, ok := err.(*domain.RateLimitError); ok {
			httpStatus = http.StatusTooManyRequests
		} else {
			httpStatus = http.StatusInternalServerError
		}
		c.JSON(httpStatus, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Notificación enviada correctamente"})
}
