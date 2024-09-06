package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error any `json:"error"`
}

type Error struct {
	Message string `json:"message"`
}

func WithInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: Error{
			Message: "что-то пошло не так, попробуйте позже.",
		},
	})
}

func WithNotFoundError(c *gin.Context, errMsg string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: Error{
			Message: errMsg,
		},
	})
}
