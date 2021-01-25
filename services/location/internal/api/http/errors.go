package httpapi

import "github.com/gin-gonic/gin"

// HTTPError model. Contains HTTP status code and a message describing the error.
type HTTPError struct {
	// HTTP status code.
	Code int `json:"code" example:"400"`
	// String describing the error that occurred.
	Message string `json:"message" example:"Bad Request"`
}

func newError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}
