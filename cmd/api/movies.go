package main

import (
	"encoding/json"
	"fmt"
	"io"
	"khomeapi/internal/data"
	"net/http"
	"time"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create a new movie ")
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	// err := json.NewDecoder(r.Body).Decode(&input)
	// to
	err := app.readJSON(w, r, &input)
	fmt.Println(err)
	if err != nil {
		app.bedRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	err = app.models.Movie.Insert(movie)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Add("Location", fmt.Sprintf("v1/movies/%d", movie.ID))
	err = app.writeJSON(w, http.StatusCreated, envelop{"movies": movie}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// fmt.Fprintln(w, "%+v\n", input)

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Runtime:   102,
		Genres:    []string{"drams", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) exampleHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Foo  string `json:"foo"`
		Test string `json:"test"`
	}
	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = json.Unmarshal(body, &input)

	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(input)

}
