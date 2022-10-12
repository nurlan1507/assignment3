package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	//routes
	router.HandlerFunc(http.MethodGet, "/", app.Home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.ViewSnippet)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.SnippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.SnippetCreatePost)

	// Create the middleware chain as normal.
	standard := alice.New(app.logRequest, secureHeader)
	return standard.Then(router)
}
