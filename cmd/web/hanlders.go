package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{"./ui/html/pages/home.tmpl", "./ui/html/partials/nav.tmpl", "./ui/html/base.tmpl"}
	ts, e := template.ParseFiles(files...)
	app.infoLog.Println("hello from handlers")
	if e != nil {
		app.serverError(w, e)
		return
	}
	err := ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
	return
}

func (app *application) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header()["X-XSS-Protection"] = []string{"1; mode=bloc"}
		w.Header().Set("Allow", "POST")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Cache-Control", "public")
		app.clientError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	w.Write([]byte("creating a new snippet"))
	return
}

func (app *application) ViewSnippet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("id")
	id, err := strconv.Atoi(query)
	if err != nil {
		app.errorLog.Println("Error id is not valid")
		http.Error(w, "Error id is not valid", http.StatusBadRequest)
		//http.NotFound(w,r)
		return
	}
	app.infoLog.Printf("Display a specific snippet with ID %d...", id)
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
