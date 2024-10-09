package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/api"
	"forum/config"
	"forum/handlers"
	"forum/utils"
)

func main() {
	if err := utils.InitServices(); err != nil {
		log.Fatal(err)
	}
	if err := utils.InitTables(); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/static/", handlers.ServeStatic)
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/api/register", api.RegisterApi)
	fmt.Printf("Server running on http://localhost%v", config.ADDRS)
	err := http.ListenAndServe(config.ADDRS, nil)
	if err != nil {
		log.Fatal(err)
	}
}
