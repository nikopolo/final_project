package server

import (
	"log"
	"net/http"

	"task_scheduler/pkg/api"
)

const DefaultPort = "7540"

func Run() error {
	api.Init()

	// файлы из директории ./web
	fs := http.FileServer(http.Dir("./web"))

	// index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if _, err := http.Dir("./web").Open(path); err != nil {
			http.ServeFile(w, r, "./web/index.html")
		} else {
			fs.ServeHTTP(w, r)
		}
	})

	log.Println("Server start at port 7540")
	return http.ListenAndServe(":7540", nil)
}
