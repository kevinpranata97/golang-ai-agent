package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"my-awesome-app/internal/config"
	"my-awesome-app/internal/database"
	"my-awesome-app/internal/handlers"
	"my-awesome-app/internal/routes"
)

func main() {
	cfg := config.Load()

	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer db.Close()

	h := handlers.New(db)

	r := gin.Default()
	routes.Setup(r, h)

	port := fmt.Sprintf(":%v", cfg.Port)
	log.Printf("Server listening on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
