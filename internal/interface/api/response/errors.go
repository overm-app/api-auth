package response

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
)

type errorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func HandleError(c *gin.Context, err error) {
	_ = c.Error(err)

	if errors.Is(err, context.Canceled) {
		c.Abort()
		return
	}

	appErr, ok := errors.AsType[*appErrors.AppError](err)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorBody{
			Code:    string(appErrors.ErrInternal),
			Message: "An unexpected error occurred",
		})
		return
	}

	body := errorBody{Code: string(appErr.Code), Message: appErr.Message}
	if appErr.HTTPStatus >= 500 {
		if appErr.Code == appErrors.ErrInternal {
			body.Message = "An unexpected error occurred"
		}
	} else {
		body.Details = appErr.Details
	}
	c.AbortWithStatusJSON(appErr.HTTPStatus, body)
}
