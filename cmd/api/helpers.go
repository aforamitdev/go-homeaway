package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type envelop map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id params")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {

	fmt.Println("new one ", dst)
	// maxBytes := 1_048_576
	// r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// dec := json.NewDecoder(r.Body)
	// dec.DisallowUnknownFields()

	// err := dec.Decode(dst)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("React Erro ")
	}

	fmt.Println(string(body[:]))
	err = json.Unmarshal(body, dst)
	fmt.Println(dst)
	fmt.Println(err, "error")
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalFieldError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		// for control maxBytesError variable
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains body-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains body-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			// if unmarshalTypeError.Field != "" {
			// 	return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			// }
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Field.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger then %d", maxBytesError.Limit)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown field %s", fieldName)

		default:
			fmt.Println("error here", err)
			return err

		}
	}

	// err = dec.Decode(&struct{}{})
	// if !errors.Is(err, io.EOF) {
	// 	return errors.New("body much not be empty")
	// }

	return nil

}

// the readString helpers

func (app *application) readString(qs url.Values, key string, defaultValue string) string {

	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}
	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")

}

func (app *application) readInt(qs url.Values, key string, defaultValue int) int {

	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)

	fmt.Println(s)
	fmt.Println(i)

	fmt.Println(key)
	fmt.Println(defaultValue)

	if err != nil {
		return defaultValue
	}
	return i
}
