package main

import (
	"log"
	"net/http"
	"os"

	apihttp "github.com/DuongThanhTin/tikfood-myself/apps/api/internal/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18081"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: apihttp.NewRouter(),
	}

	log.Printf("tikfood api listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("api server stopped: %v", err)
	}
}
