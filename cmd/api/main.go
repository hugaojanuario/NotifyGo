package main

import (
	"log"

	"github.com/hugaojanuario/NotifyGo/internal/handler"
	"github.com/hugaojanuario/NotifyGo/internal/repository"
	"github.com/hugaojanuario/NotifyGo/internal/router"
	"github.com/hugaojanuario/NotifyGo/internal/service"
	"github.com/hugaojanuario/NotifyGo/pkg/config"
	"github.com/hugaojanuario/NotifyGo/pkg/database"
)

func main() {
	cfg := config.LoadDotEnv()

	db, err := database.DBConn(cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close()

	r := repository.NewUserRepository(db)
	s := service.NewUserService(r)
	h := handler.NewUserHandler(s)

	router := router.SetupRouter(h)
	router.Run(":9292")
}
