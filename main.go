package main

import (
	"log"
	"music-library/config"
	"music-library/controllers"
	"music-library/migrations"
	"music-library/routes"
	"net/http"

	"github.com/go-chi/chi"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func startServer(router *chi.Mux, port string) {
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("server starting on port %v", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("couldn't start server %v", err)
	}
}

func main() {
	config.LoadEnv()

	config := config.LoadConfig()

	db, err := gorm.Open(postgres.Open(config.DB_URL), &gorm.Config{})
	if err != nil {
		log.Fatalf("couldn't connect to database: %v", err)
	}

	migrations.Migrate(db)

	controllers.SetDB(db)

	router := routes.SetupRoutes()

	log.Println("database connected successfully")
	startServer(router, config.PORT)
}