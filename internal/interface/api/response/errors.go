package response

import (
	"errors"

	"github.com/gin-gonic/gin"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
)

func HandleError(c *gin.Context, err error) {
    _ = c.Error(err)
		if appErr, ok := errors.AsType[*appErrors.AppError](err); ok {
		c.AbortWithStatusJSON(appErr.HTTPStatus, gin.H{
			"code":    string(appErr.Code),
			"message": appErr.Message,
		})
		return
	}
	c.AbortWithStatusJSON(500, gin.H{
		"code":    string(appErrors.ErrInternal),
		"message": "An unexpected error occurred",
	})
}
