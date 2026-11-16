package response

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
)

func BindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindBodyWithJSON(dst); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			details := map[string]any{}
			for _, fe := range ve {
				details[strings.ToLower(fe.Field())] = fe.Tag()
			}
			HandleError(c, appErrors.Validation(appErrors.ErrValidation, "Request body is invalid").WithDetails(details))
			return false
		}
		HandleError(c, appErrors.Validation(appErrors.ErrBadRequest, "Malformed request body"))
		return false
	}
	return true
}
