package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// Server error helper writes error message and stack trace
// and also sends generic server error (500) via response to user
func (app *application) serverError(w http.ResponseWriter, err error) {
	// debug.Stack is stack trace of current goroutine
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// client error helper sends status code to user
// code 400
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound helper wrapper to clientError 404 not found
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// retrieve corresponding template set cached based on page
	// raise error if page does not exist
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialise buffer to simulate response
	buf := new(bytes.Buffer)

	// exec template
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// write status out to response
	w.WriteHeader(status)

	// write to user response
	buf.WriteTo(w)
}

// initialises current year automatically
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// automate display of flash messages
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call decoder
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// check if invalid
		var invalidDecoderError *form.InvalidDecoderError

		// panic if invalid target destination
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// all other errors
		return err
	}

	return nil
}
