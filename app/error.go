package app

import (
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int
	Name       string
	Message    string
}

type ErrorType int

const (
	NetworkError ErrorType = iota
	AuthError
	ParsingError
	InvalidCourseIdError
	DatabaseError
	CryptographyError
	UnknownError
)

//TODO: Add ability to add detailed error message from throwing call

func CreateError(errorType ErrorType) Error {
	switch errorType {
	case NetworkError:
		return Error{
			StatusCode: http.StatusBadGateway,
			Name:       "Network Error",
			Message:    "A network error has occurred while communicating with TeachAssist.",
		}
	case AuthError:
		return Error{
			StatusCode: http.StatusUnauthorized,
			Name:       "Authentication Error",
			Message:    "An authentication error has occurred. Invalid username or password.",
		}
	case ParsingError:
		return Error{
			StatusCode: http.StatusInternalServerError,
			Name:       "Parsing Error",
			Message:    "An error has occurred while parsing the response from TeachAssist.",
		}
	case InvalidCourseIdError:
		return Error{
			StatusCode: http.StatusBadRequest,
			Name:       "Invalid Course Id Error",
			Message:    "An invalid course id was given.",
		}
	case DatabaseError:
		return Error{
			StatusCode: http.StatusInternalServerError,
			Name:       "Database Error",
			Message:    "An error has occurred while accessing the database.",
		}
	case CryptographyError:
		return Error{
			StatusCode: http.StatusInternalServerError,
			Name:       "Cryptography Error",
			Message:    "An error has occurred while performing a cryptographic function.",
		}
	case UnknownError:
		return Error{
			StatusCode: http.StatusInternalServerError,
			Name:       "Unknown Error",
			Message:    "An unknown error has occurred.",
		}
	default:
		return Error{
			StatusCode: http.StatusInternalServerError,
			Name:       "Unknown Error",
			Message:    "An unknown error has occurred.",
		}
	}
}

func (e Error) String() string {
	return fmt.Sprintf("(%v) %v: %v", e.StatusCode, e.Name, e.Message)
}

func (e Error) Error() string {
	return e.String()
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (e Error) ErrorResponse() ErrorResponse {
	return ErrorResponse{
		Error:   e.Name,
		Message: e.Message,
	}
}
