package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ysugimoto/doorkeeper/handler"
)

func main() {
	port := "9000"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	if err := http.ListenAndServe(":"+port, handler.WebhookHandler("/webhook/")); err != nil {
		log.Fatalln(err)
	}
}
