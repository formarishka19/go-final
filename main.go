package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	"go-final/pkg/db"
	"go-final/pkg/server"

	_ "modernc.org/sqlite"
)

const dbFile = "scheduler.db"

func main() {
	install := false
	_, err := os.Stat(dbFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			install = true
		} else {
			log.Fatal("error accessing db file, possible no access %w", err)
			return
		}
	}
	db.DB, err = sql.Open("sqlite", dbFile)
	defer db.DB.Close()
	if err != nil {
		log.Fatal("error acces to database %w", err)
		return
	}
	if install {
		err = db.Init(dbFile)
		if err != nil {
			log.Fatal("error generating db schema %w", err)
			return
		}
	}

	httpLogger := new(log.Logger)
	httpServer := server.NewHttpServer(httpLogger)
	err = httpServer.Srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		httpLogger.Fatalf("Ошибка сервера: %v", err)
	}
}
