package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/recipe/view", app.recipeView)
	mux.HandleFunc("/recipe/create", app.recipeCreate)

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
