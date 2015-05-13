package infrared

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

func NewAPIError() APIError {
	return APIError{}
}

type APIError struct {
	Code    int    `json:"code"`
	Title   string `json:"error"`
	Message string `json:"message"`
}

func (e APIError) code(code int) APIError {
	e.Code = code
	return e
}

func (e APIError) title(title string) APIError {
	e.Title = title
	return e
}

func (e APIError) message(message string) APIError {
	e.Message = message
	return e
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Message)
}

func (e APIError) FromQueryError(err error) APIError {
	switch err.Error() {
	case "bad request":
		e.Code = 400
		e.Title = "Bad Request"
		e.Message = "The request you made contained invalid data, please check it and try again."
	case "not found":
		e.Code = 404
		e.Title = "Not Found"
		e.Message = "We could not find an entry which matched the criteria you specified. Please check them and try again."

	default:
		log.Printf("Unknown error %s", err)
		e.Code = 500
		e.Title = err.Error()
		e.Message = "An unknown database error occurred when processing your request. Please try again at a later stage."

	}

	return e
}

func (e APIError) ToResponse(res rest.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")

	switch e.Code {
	case 400:
		res.WriteHeader(http.StatusBadRequest)
	case 401:
		res.WriteHeader(http.StatusUnauthorized)
	case 403:
		res.WriteHeader(http.StatusForbidden)
	case 404:
		res.WriteHeader(http.StatusNotFound)
	case 405:
		res.WriteHeader(http.StatusMethodNotAllowed)
	case 409:
		res.WriteHeader(http.StatusConflict)
	case 500:
		fallthrough
	default:
		res.WriteHeader(http.StatusInternalServerError)
	}

	res.WriteJson(e)
}
