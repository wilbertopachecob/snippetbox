package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"wilbertopachecob/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//Check if the current request URL path exactly matches "/". If it doesn't
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the hand
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	// Working Directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		app.errorlog.Fatal(err)
	}

	files := []string{
		dir + "/ui/html/home.page.tmpl",
		dir + "/ui/html/footer.partial.tmpl",
		dir + "/ui/html/base.layout.tmpl"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.infolog.Println(err.Error())
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.infolog.Println(err.Error())
		app.serverError(w, err)
		return
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || ID <= 0 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(ID)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Fprintf(w, "%v", s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(405)
		// w.Write([]byte("Method not allowed"))
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
	expires := "7"

	ID, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", ID), http.StatusSeeOther)
}
