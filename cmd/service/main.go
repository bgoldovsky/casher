package main

import (
	"database/sql"
	"fmt"
	"net/http"

	operationsRepo "github.com/bgoldovsky/casher/app/repositories/operations"
	usersRepo "github.com/bgoldovsky/casher/app/repositories/users"
	"github.com/bgoldovsky/casher/app/services/operations"
	"github.com/bgoldovsky/casher/app/services/users"
	"github.com/bgoldovsky/casher/config"
	"github.com/bgoldovsky/casher/handlers"
	_ "github.com/lib/pq"
)

// TODO: Сделать все красивым HTML+CSS

// Запускаем сервер
func handleRequest(handler *handlers.PageHandler, port string) {
	addr := fmt.Sprintf(":%s", port)
	if err := http.ListenAndServe(addr, handler.Router()); err != nil {
		panic(err)
	}
}

func main() {
	// Инициализируем БД
	connString := config.ConnectionString()
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Repositories
	operationsRepository := operationsRepo.New(db)
	usersRepository := usersRepo.New(db)

	// Services
	usersSrv := users.New(usersRepository, operationsRepository)
	operationsSrv := operations.New(operationsRepository)

	// Handlers
	htmlHandler := handlers.New(usersSrv, operationsSrv)

	// Запуск сервера
	port := config.Port()
	handleRequest(htmlHandler, port)
}
