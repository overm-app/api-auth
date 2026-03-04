package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
)

func HandleError(c *gin.Context, sugar *zap.SugaredLogger, err error, requestID string) {
    var appErr *appErrors.AppError

    if errors.As(err, &appErr) {
        if appErr.Err != nil {
            sugar.Errorw("Request failed",
                "request_id", requestID,
                "code",       appErr.Code,
                "cause",      appErr.Err.Error(),
            )
        }
        c.AbortWithStatusJSON(appErr.HTTPStatus, gin.H{
            "code":    string(appErr.Code),
            "message": appErr.Message,
        })
        return
    }

    sugar.Errorw("Unexpected error",
        "request_id", requestID,
        "error", err.Error(),
    )
    c.AbortWithStatusJSON(500, gin.H{
        "code":    string(appErrors.ErrInternal),
        "message": "An unexpected error occurred",
    })
}