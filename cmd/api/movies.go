package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"khomeapi/internal/data"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	// err := json.NewDecoder(r.Body).Decode(&input)
	// to
	// body, err := io.ReadAll(r.Body)
	// fmt.Println(body)
	fmt.Println("before reactJSON")
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.bedRequestResponse(w, r, err)
		return
	}

	fmt.Println(input.Title, input.Genres, input.Year)
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	fmt.Println("Text one one ", movie)

	err = app.models.Movie.Insert(movie)
	fmt.Println(err)
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
	fmt.Println(id)

	movie, err := app.models.Movie.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movie.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
			return

		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return

		default:
			app.serverErrorResponse(w, r, err)
			return
		}
		return
	}
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.bedRequestResponse(w, r, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = int32(input.Runtime)
	movie.Genres = input.Genres

	err = app.models.Movie.Update(movie)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movie.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelop{"message": "movies successfully deleted"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title  string
		Genres []string
		data.Filter
	}

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	// sort
	input.Sort = app.readString(qs, "sort", "id")

	movies, err := app.models.Movie.GetAll(input.Title, input.Genres, input.Filter)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"movies": movies}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
