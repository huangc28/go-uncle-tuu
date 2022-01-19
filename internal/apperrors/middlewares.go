package apperrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()

		if err != nil {
			srcErr := err.Err
			var parsedErr *Error

			switch srcErr := srcErr.(type) {
			case *Error:
				parsedErr = srcErr
			default:
				parsedErr = &Error{
					ErrCode: UnknownErrorToApplication,
					Err:     srcErr.Error(),
				}

				c.Writer.WriteHeader(http.StatusInternalServerError)
			}

			c.JSON(c.Writer.Status(), parsedErr)

			return
		}
	}
}
