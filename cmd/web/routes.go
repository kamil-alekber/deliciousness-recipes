package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)

	mux.HandleFunc("/login", app.login)

	mux.HandleFunc("/recipe/view", app.recipeView)
	mux.HandleFunc("/recipe/create", app.recipeCreate)

	mux.HandleFunc("/login/google", app.loginGoogle)
	mux.HandleFunc("/login/google/redirect", app.loginGoogleRedirect)

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
