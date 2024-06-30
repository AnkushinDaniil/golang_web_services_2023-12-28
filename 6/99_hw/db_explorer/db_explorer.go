package main

import (
	"database/sql"
	"net/http"

	"db_explorer/handler"
	"db_explorer/repository"
	"db_explorer/service"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	repo, err := repository.NewRepository(db)
	if err != nil {
		return nil, err
	}
	services := service.NewService(repo.Item, repo.Table)
	handlers := handler.NewHandler(services.Item, services.Table)
	return handlers, nil
}
