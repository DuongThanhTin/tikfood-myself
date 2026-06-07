package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Data  any            `json:"data"`
	Error *errorResponse `json:"error,omitempty"`
}

type errorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func respondWithData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response{Data: data})
}

func respondWithError(c *gin.Context, status int, code string, message string) {
	c.JSON(status, response{
		Data: nil,
		Error: &errorResponse{
			Code:    code,
			Message: message,
		},
	})
}

func respondWithBadRequest(c *gin.Context, err *errorResponse) {
	c.JSON(http.StatusBadRequest, response{Data: nil, Error: err})
}

func respondWithInternalServerError(c *gin.Context, message string) {
	respondWithError(c, http.StatusInternalServerError, ErrorCodeInternal, message)
}
