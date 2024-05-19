package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Define new servemux (the map between URL patterns and handlers),
// then register the home function as the handler for '/'
func (app *application) routes() http.Handler {
	// NOTE:
	// Could actually avoid declaring this and just use an http.HandleFunc, http.ListenAndServe
	// right away. http will set a global DefaultServeMux under the hood.
	// WARNING: not a good a idea for production applications, as any 3rd party library
	// can access the global mux and add routes to it. If a 3rd party is compromised
	// this would make the Go application vulnerable.
	// mux := http.NewServeMux()

	// use httprouter
	router := httprouter.New()

	// hook for custom exceptions
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// NOTE: create a File Server to serve static files
	// here we want it to be a subtree path, so add a trailing /
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)


	// setup mux to handle all requests to 'static'
	// we strip prefix before request reaches the file server
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// '/' is a subtree path, i.e. it will match with any path starting with '/'
	// and trailing with anything or nothing at all.
	// mux.HandleFunc("/", app.home)

	// Fixed paths, i.e. only match the URL string exactly.
	// NOTE: longer patterns are always matched first
	// mux.HandleFunc("/snippet/view", app.snippetView)
	// mux.HandleFunc("/snippet/create", app.snippetCreate)
	// NOTE: URL request patterns are sanitised automatically.
	// Any . or .. or // will be removed and the user redirected to
	// the equivalent clean URL
	// ALSO WORKS WITH PATH TYPES: if a user forgot a trailing slash for
	// a subtree path, this will also be sanitised

	// create standard chain of middleware (default)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// finally serve http map (mux)
	return standard.Then(router)
}
