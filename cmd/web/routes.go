package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Define new servemux (the map between URL patterns and handlers),
// then register the home function as the handler for '/'
func (app *application) routes() http.Handler {
	// use httprouter
	router := httprouter.New()

	// hook for custom exceptions
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// NOTE: create a File Server to serve static files
	// here we want it to be a subtree path, so add a trailing /
	// static files do not need a stateful session (so no sessionManager)
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// new middleware containing session data
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// create standard chain of middleware (default)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// finally serve http map (mux)
	return standard.Then(router)
}
