package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/crazyfacka/yanbapp-api/config"
	"github.com/crazyfacka/yanbapp-api/repositories/cache"
	"github.com/crazyfacka/yanbapp-api/repositories/db"
	"github.com/crazyfacka/yanbapp-api/services/api"
)

var (
	dbRepository    *db.DB
	cacheRepository *cache.Redis
)

func listenInterrupt(done chan bool) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("yanbapp-api is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	api.Stop(ctx)
	dbRepository.Close()
	cacheRepository.Close()

	close(done)
}

func main() {
	done := make(chan bool)
	go listenInterrupt(done)

	config.LoadConfiguration()

	dbconf := config.DB()
	dbRepository = db.NewDB(dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Schema)

	redisconf := config.Cache()
	cacheRepository = cache.NewRedis(redisconf.Host, redisconf.Port, redisconf.Database)

	api.Start(config.API().Port, dbRepository, cacheRepository)

	<-done
	log.Println("all stopped")
}
