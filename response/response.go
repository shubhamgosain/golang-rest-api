package response

import (
	"net/http"
	"github.com/go-chi/render"
)

type response struct {
	HTTPStatusCode int         `json:"http_status_code"` // http response status code
	StatusText     string      `json:"status"`           // user-level status message
	ErrorText      string      `json:"error"`            // application-level error message, for debugging
	Data           interface{} `json:"data"`             // application-level data
}

func (e *response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// Bad Request
func ErrBadRequest(err error) render.Renderer {
	return &response{
		HTTPStatusCode: 400,
		StatusText:     "Bad request",
		ErrorText:      err.Error(),
	}
}

func errInternalServer(err error) render.Renderer {
	return &response{
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}

func SuccessResponse(data interface{}) render.Renderer {
	return &response{
		HTTPStatusCode: 200,
		StatusText:     "OK",
		Data:           data,
	}
}