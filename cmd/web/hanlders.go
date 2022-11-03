package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/nurlan1507/internal/models"
	"net/http"
	"strconv"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {

	snippets, snErr := app.snippets.Latest()
	if snErr != nil {
		app.errorLog.Println("SMTH HAPPEN TO SNIPPETS FETCH")
		app.serverError(w, snErr)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
	return
}

func (app *application) SnippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	app.render(w, 200, "create.tmpl", data)
}

func (app *application) SnippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	w.Header()["X-XSS-Protection"] = []string{"1; mode=bloc"}
	//	w.Header().Set("Allow", "POST")
	//	w.Header().Set("Content-Type", "application/json")
	//	w.Header().Add("Cache-Control", "public")
	//	_, err := app.snippets.Insert("nurik", "contentLOL", 24)
	//	if err != nil {
	//		app.errorLog.Println("some error inserting to db")
	//		return
	//	}
	//	app.clientError(w, http.StatusMethodNotAllowed, "Method not allowed")
	//	return
	//}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	newSnippet, err := app.snippets.Insert(title, content, 5)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/:%d", newSnippet.ID), http.StatusSeeOther)
	return
}

func (app *application) ViewSnippet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorLog.Println("Error id is not valid")
		http.Error(w, "Error id is not valid", http.StatusBadRequest)
		//http.NotFound(w,r)
		return
	}
	fmt.Println(id)
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	data := app.newTemplateData(r) //чтобы кидать больше объектов в шаблон
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
	app.infoLog.Printf("Display a specific snippet with ID %d...", id)
}
