package main

import (
	"log"

	"task_scheduler/pkg/db"
	"task_scheduler/pkg/server"
)

func main() {
	// проверка и запуск базы
	err := db.Init("scheduler.db")
	// закрытие подключения к БД
	defer db.CloseDB()
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	// запуск Server
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
