package routes

import (
	"music-library/controllers"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/song/list", controllers.GetSongs)
	router.Get("/song/{id}/text", controllers.GetSongText)
	router.Delete("/song/{id}/delete", controllers.DeleteSong)
	router.Put("/song/{id}/edit", controllers.EditSong)
	router.Post("/song/create", controllers.AddSong)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return router
}