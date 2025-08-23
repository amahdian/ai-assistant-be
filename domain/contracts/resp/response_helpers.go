package resp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/amahdian/ai-assistant-be/domain/model/common"
	"github.com/amahdian/ai-assistant-be/global/errs"
	"github.com/amahdian/ai-assistant-be/pkg/logger"
	"github.com/amahdian/ai-assistant-be/pkg/msg"
	"github.com/amahdian/ai-assistant-be/server/binding"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

func AbortWithError(ctx *gin.Context, err error) {
	var ve validator.ValidationErrors
	var customErr *errs.Error
	switch {
	case errors.As(err, &ve):
		mc := binding.MapValidationErrorsToMessageContainer(ve)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(mc))
		return
	case errors.As(err, &customErr):
		ctx.AbortWithStatusJSON(customErr.Code.HttpStatus(), NewErrorResponse(err))
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}
}

func AbortWithStatus(ctx *gin.Context, code int, data interface{}) {
	switch v := data.(type) {
	case error:
		{
			var ve validator.ValidationErrors
			switch {
			case errors.As(v, &ve):
				mc := binding.MapValidationErrorsToMessageContainer(ve)
				ctx.AbortWithStatusJSON(code, NewErrorResponse(mc))
				return
			default:
				ctx.AbortWithStatusJSON(code, NewErrorResponse(v))
				return
			}
		}
	case *msg.MessageContainer:
		{
			ctx.AbortWithStatusJSON(code, NewErrorResponse(v))
			return
		}
	default:
		logger.Errorf("could not recognize the error type for proper error handling: %v", data)
	}
}

func Success(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, NewResponse(true))
}

func Ok(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, NewResponse(data))
}

func OkWithMessage(ctx *gin.Context, data interface{}, messages *msg.MessageContainer) {
	ctx.JSON(http.StatusOK, NewResponseWithMessage(data, messages))
}

func PaginatedOk[T any](ctx *gin.Context, data []T, pagination *common.Pagination) {
	ctx.JSON(http.StatusOK, NewPaginatedResponse(data, pagination))
}

func Created(ctx *gin.Context, data ...interface{}) {
	if len(data) > 0 {
		ctx.JSON(http.StatusCreated, NewResponse(data[0]))
	} else {
		ctx.JSON(http.StatusCreated, NewResponse(true))
	}
}

func NoContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func CreatedWithMessage(ctx *gin.Context, data interface{}, messages *msg.MessageContainer) {
	ctx.JSON(http.StatusCreated, NewResponseWithMessage(data, messages))
}

func WriteSheetBytes(ctx *gin.Context, data []byte, fileName string) {
	setSheetHeaders(ctx, fileName)
	_, _ = ctx.Writer.Write(data) // ignore error
}

func WriteJsonFileBytes(ctx *gin.Context, data []byte, fileName string) {
	fileName = strings.TrimSuffix(fileName, ".json") // remove suffix ".json" if present
	fileName = fmt.Sprintf("%s.json", fileName)      // add suffix ".json"
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q; filename*=utf-8''%q", fileName, fileName))
	_, _ = ctx.Writer.Write(data) // ignore error
}

func setSheetHeaders(ctx *gin.Context, fileName string) {
	fileName = strings.TrimSuffix(fileName, ".xlsx") // remove suffix ".xlsx" if present
	fileName = fmt.Sprintf("%s.xlsx", fileName)      // add suffix ".xlsx"
	ctx.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q; filename*=utf-8''%q", fileName, fileName))
}

func Redirect(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusSeeOther, location)
}
