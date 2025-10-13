package main

import (
	"log"

	"task_scheduler/pkg/db"
	"task_scheduler/pkg/server"
)

func main() {
	// проверка и запуск базы
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	// // файлы из директории ./web
	// fs := http.FileServer(http.Dir("./web"))

	// // Оборачиваем FileServer, чтобы перенаправлять все не найденные пути на index.html
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	path := r.URL.Path

	// 	// Проверяем, существует ли запрашиваемый файл
	// 	if _, err := http.Dir("./web").Open(path); err != nil {
	// 		// Если файл не найден, возвращаем index.html (для SPA)
	// 		http.ServeFile(w, r, "./web/index.html")
	// 	} else {
	// 		// Иначе отдаём запрошенный файл
	// 		fs.ServeHTTP(w, r)
	// 	}
	// })

	// http.ListenAndServe(":7540", nil)

	// запуск Server
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
