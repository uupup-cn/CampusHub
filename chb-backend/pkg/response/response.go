package response

import (
	"net/http"

	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedData struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func SuccessPaginated(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PaginatedData{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func Error(c *gin.Context, ec errcode.ErrorCode) {
	c.JSON(ec.HTTP, Response{
		Code:    ec.Code,
		Message: ec.Message,
	})
}

func ErrorWithMessage(c *gin.Context, ec errcode.ErrorCode, msg string) {
	c.JSON(ec.HTTP, Response{
		Code:    ec.Code,
		Message: msg,
	})
}

func ServerError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.JSON(http.StatusInternalServerError, Response{
		Code:    5001,
		Message: "系统内部错误",
	})
}
