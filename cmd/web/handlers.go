package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"wilbertopachecob/snippetbox/pkg/forms"
	"wilbertopachecob/snippetbox/pkg/models"

	"github.com/justinas/nosurf"
)

func getExecutablePath() string {
	// Working Directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "."
	}
	return dir
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{CurrentYear: time.Now().Year()}
		return td
	}
	td.CurrentYear = time.Now().Year()
	// Use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it
	// acts like a one-time fetch. If there is no matching key in the session
	// data this will return the empty string.
	td.Flash = app.session.PopString(r, "flash")
	td.Authenticated = app.authenticatedUser(r)
	td.CSRFToken = nosurf.Token(r)
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, page string, data *templateData) {
	// Retrieve the appropriate template set from the cache based on the page n
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.
	ts, ok := app.templateCache[page]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", page))
		return
	}
	// Initialize a new buffer.
	buf := new(bytes.Buffer)
	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError helper and
	// return.
	err := ts.Execute(buf, app.addDefaultData(data, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//Check if the current request URL path exactly matches "/". If it doesn't
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the hand
	// would keep executing and also write the "Hello from SnippetBox" message.
	// Because Pat matches the "/" path exactly, we can now remove the manual c
	// of r.URL.Path != "/" from this handler.
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: snippets}
	// Working Directory
	// dir := getExecutablePath()

	// files := []string{
	// 	dir + "/ui/html/home.page.tmpl",
	// 	dir + "/ui/html/footer.partial.tmpl",
	// 	dir + "/ui/html/base.layout.tmpl"}
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	ID, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

	data := &templateData{Snippet: s}

	app.render(w, r, "show.page.tmpl", data)

	// dir := getExecutablePath()
	// files := []string{
	// 	dir + "/ui/html/show.page.tmpl",
	// 	dir + "/ui/html/base.layout.tmpl",
	// 	dir + "/ui/html/footer.partial.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	// app.render("create.page.tmpl", w, nil, r)
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// The check of r.Method != "POST" is now superfluous and can be removed.
	// if r.Method != "POST" {
	// 	w.Header().Set("Allow", "POST")
	// 	// w.WriteHeader(405)
	// 	// w.Write([]byte("Method not allowed"))
	// 	//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError helper to
	// a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// expires := r.PostForm.Get("expires")

	// errors := make(map[string]string)

	// if strings.TrimSpace(title) == "" {
	// 	errors["title"] = "The field can not be blank"
	// } else if utf8.RuneCountInString(title) > 100 {
	// 	errors["title"] = "The field can not contain more than 100 characters"
	// }

	// if strings.TrimSpace(content) == "" {
	// 	errors["content"] = "The field can not be blank"
	// }

	// if strings.TrimSpace(expires) == "" {
	// 	errors["expires"] = "The field can not be blank"
	// } else if expires != "1" && expires != "7" && expires != "365" {
	// 	errors["expires"] = "This field is invalid"
	// }
	// // If there are any errors, dump them in a plain text HTTP response and ret
	// // from the handler.
	// if len(errors) > 0 {
	// 	data := &templateData{
	// 		FormData:   r.PostForm,
	// 		FormErrors: errors,
	// 	}
	// 	app.render("create.page.tmpl", w, data, r)
	// 	return
	// }

	form := forms.New(r.PostForm)

	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermitedValues("expires", "1", "7", "365")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	ID, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Your snippet was saved
	// successfully!") and the corresponding key ("flash") to the session
	// data. Note that if there's no existing session for the current user
	// (or their session has expired) then a new, empty, session for them
	// will automatically be created by the session middleware.
	app.session.Put(r, "flash", "The Snippet was created successfuly")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", ID), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
	}

	f := forms.New(r.PostForm)
	f.Required("name", "email", "password")
	f.MinLength("password", 10)
	f.MatchesPattern("email", forms.EmailRX)

	if !f.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: f,
		})
		return
	}

	err = app.users.Insert(f.Get("name"), f.Get("email"), f.Get("password"))
	if err != nil {
		if err == models.ErrDuplicateEmail {
			f.Errors.Add("email", "This email already exist on the DB")
			app.render(w, r, "signup.page.tmpl", &templateData{
				Form: f,
			})
			return
		}
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	f := forms.New(r.PostForm)
	f.Required("email", "password")
	f.MatchesPattern("email", forms.EmailRX)

	if !f.Valid() {
		app.render(w, r, "login.page.tmpl", &templateData{
			Form: f,
		})
		return
	}

	id, err := app.users.Authenticate(f.Get("email"), f.Get("password"))
	if err != nil {
		if err == models.ErrNoRecord {
			f.Errors.Add("generic", "There is no user on the DB with the provided email")
			app.render(w, r, "login.page.tmpl", &templateData{Form: f})
			return
		} else if err == models.ErrInvalidCredentials {
			f.Errors.Add("generic", "Please verify the provided password")
			app.render(w, r, "login.page.tmpl", &templateData{Form: f})
			return
		}
		app.serverError(w, err)
		return
	}
	// Add the ID of the current user to the session, so that they are now 'logg
	// in'.
	app.session.Put(r, "userID", id)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "You have been logout successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
