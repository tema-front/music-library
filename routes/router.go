package routes

import (
	"music-library/controllers"

	"github.com/go-chi/chi"
)

func SetupRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/song/list", controllers.GetSongs)
	router.Get("/song/{id}/text", controllers.GetSongText)
	router.Delete("/song/{id}/delete", controllers.DeleteSong)
	router.Put("/song/{id}/edit", controllers.EditSong)
	router.Post("/song/create", controllers.AddSong)

	return router
}