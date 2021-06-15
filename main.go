package main

import (
	"github.com/cmelgarejo/minesweeper-svc/database"
	"github.com/cmelgarejo/minesweeper-svc/utils/config"
	"github.com/cmelgarejo/minesweeper-svc/utils/logger"
	server "github.com/cmelgarejo/minesweeper-svc/web"
	"github.com/loopcontext/msgcat"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg.Debug)
	catalog, err := msgcat.NewMessageCatalog(cfg.MessageCatalog)
	if err != nil {
		log.SendFatal(err)
	}
	db, err := database.NewDB(cfg.DB, log)
	if err != nil {
		log.SendFatal(err)
	}
	appsrv, err := server.InitFiberServer(cfg, log, &catalog, db)
	if err != nil {
		log.SendFatal(err)
	}
	log.SendFatal(appsrv.Listen(cfg.Server.BuildServerAddr()))
}
