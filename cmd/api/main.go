package main

import (
	"log"

	"github.com/hugaojanuario/NotifyGo/internal/server"
	"github.com/hugaojanuario/NotifyGo/internal/user"
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

	r := user.NewUserRepository(db)
	s := user.NewUserService(r)
	h := user.NewUserHandler(s)

	router := server.SetupRouter(h)
	router.Run(":9292")
}
