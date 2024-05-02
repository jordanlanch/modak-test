package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Application struct {
	Env *Env
	Rdb *redis.Client
}

func App(envfile string) *Application {
	app := &Application{
		Env: NewEnv(envfile),
	}

	var err error
	for i := 1; i <= 3; i++ {
		// Initialize Redis client
		app.Rdb = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", app.Env.RedisHost, app.Env.RedisPort),
			DB:   0, // default database
		})

		// Context with timeout for connecting to Redis
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Check Redis connection
		if _, err = app.Rdb.Ping(ctx).Result(); err == nil {
			break // La conexión fue exitosa, salir del bucle de reintento
		}

		time.Sleep(1 * time.Second) // Esperar antes de intentar nuevamente
	}

	// Verificar si se pudo conectar a Redis
	if err != nil {
		log.Fatalf("Unable to connect to Redis after retries: %s", err)
	}

	return app
}
