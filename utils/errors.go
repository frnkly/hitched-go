package utils

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrorResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render ...
func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)

	return nil
}

// BadRequestError ...
func BadRequestError(err error) render.Renderer {
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Bad request.",
		ErrorText:      err.Error(),
	}
}

// NotFoundError ...
func NotFoundError(err error) render.Renderer {
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     "Not found.",
		ErrorText:      err.Error(),
	}
}

// RenderingError ...
func RenderingError(err error) render.Renderer {
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
