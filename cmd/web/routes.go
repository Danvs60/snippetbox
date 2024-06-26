package main

import (
	"net/http"

	"github.com/Danvs60/snippetbox/ui"

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
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// Ping route for testing
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// unprotected application routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	// user routes
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// PROTECTED user routes
	// no need to attach 'noSurf' as it appends on dynamic Handler
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// create standard chain of middleware (default)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// finally serve http map (mux)
	return standard.Then(router)
}
