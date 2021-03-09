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

	log.Printf("Server starts on :%s", port)
	if err := http.ListenAndServe(
		":"+port,
		handler.WebhookHandler("/webhook", nil),
	); err != nil {
		log.Fatalln(err)
	}
}
