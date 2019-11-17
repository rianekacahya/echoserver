package response

import (
	"github.com/labstack/echo"
	"github.com/rianekacahya/errors"
	"net/http"
)

type response struct {
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Error(c echo.Context, err error) error {
	var(
		status int
		response = new(response)
	)

	if err != nil {
		// Mapping Status
		errorStatus := errors.GetStatus(err)
		errorMessage := errors.GetError(err)
		switch errorStatus {
		case errors.GENERIC:
			status = http.StatusInternalServerError
		case errors.FORBIDDEN:
			status = http.StatusForbidden
		case errors.BADREQUEST:
			status = http.StatusBadRequest
		case errors.NOTFOUND:
			status = http.StatusNotFound
		case errors.UNAUTHORIZED:
			status = http.StatusUnauthorized
		default:
			status = http.StatusInternalServerError
			response.Message = err.Error()
		}

		if errorStatus != errors.NOTYPE {
			switch errorMessage.(type) {
			case errors.Custom:
				response.Message = errorMessage.Error()
			default:
				response.Message = errorMessage
			}
		}
	}

	return c.JSON(status, response)
}

func Render(c echo.Context, status int, data interface{}) error {
	var response = new(response)

	response.Message = "success"
	response.Data = data

	return c.JSON(status, response)
}