package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Danvs60/snippetbox/internal/models"
	"github.com/Danvs60/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// Define home handler function
// Write a byte slice containing "Hello from Snippetbox" as the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if request URL exactly matches '/'
	// If not, send the user to a 404 page.
	// Important to also return from the handler to avoid
	// writing the "Hello from Snippetbox" response.

	// WARNING: NOT NEEDED ANYMORE AS httprouter matches '/' exactly
	/* if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	} */

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)

	// =======================
	// OLD WAY TO PARSE TEMPLATE
	/* files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}
	// Use template all template files
	// Use http.Error to log any error under res code 500
	// NOTE: files is used as a variadic parameter
	// the function ParseFiles is defined (files ...str) -> meaning it can take in a variable size slice as input
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create templateData and add snippets
	data := &templateData{
		Snippets: snippets,
	}

	// Write content of base template in response body
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	} */
	// =======================
}

// Define snippetView handler function
// Writes a byte string of a snippet's content
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// get params (httprouter)
	params := httprouter.ParamsFromContext(r.Context())
	// Extract value of id parameter
	// id_param := r.URL.Query().Get("id")
	id_param := params.ByName("id")

	// Validate it is an integer
	id, err := strconv.Atoi(id_param)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
	// OLD WAY TO PARSE TEMPLATE
	// ===================================
	/* // Initialise a slice containing paths to view.tmpl
	// and other layouts
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	// Parse templates
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// struct holding the snippet data
	data := &templateData{
		Snippet: snippet,
	}

	// execute templates
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	} */
	// ===================================
}

// form struct to represent data and validation
// implemented with decoder, which extracts values from HTML form
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // - tells decoder to ignore a field during decoding
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

// Define a snippetCreate handler function
// Creates a new snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Check if request is POST
	// WARNING: step not needed, httprouter does it automatically
	/* if r.Method != http.MethodPost {
	// let user it is not allowed to use this method
	w.Header().Set("Allow", http.MethodPost) // WARNING: you can only set response headers before the first call to Write

	/* w.WriteHeader(405) // WARNING: can only call once per response
	// NOTE: if WriteHeader is not called, first time we 'Write' it will auto call to res = 200
	w.Write([]byte("Method not allowed")) */

	// This is better done with http Error
	// app.clientError(w, http.StatusMethodNotAllowed)
	// return

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}


	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// re-display create.tmpl if there are any errors to display
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// Pass snippet data to connection pool for insert
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
