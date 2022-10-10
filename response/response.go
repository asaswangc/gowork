package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PageResult struct {
	Rows     interface{} `json:"rows"`
	Code     int         `json:"code"`
	Total    int64       `json:"total"`
	PageNum  int         `json:"page_num"`
	PageSize int         `json:"page_size"`
}

type JsonResult struct {
	Rows    interface{} `json:"rows"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

func NewJsonResult(code int, message interface{}, rows interface{}) *JsonResult {
	if rows == nil {
		rows = []struct{}{}
	}
	return &JsonResult{
		Code:    code,
		Rows:    rows,
		Message: message,
	}
}

func NewPageResult(code int, total int64, pageNum int, pageSize int, rows interface{}) *PageResult {
	if rows == nil {
		rows = []struct{}{}
	}
	return &PageResult{
		Code:     code,
		Rows:     rows,
		Total:    total,
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}

type output func(ctx *gin.Context, v interface{})
type ResultFunc func(result interface{}) func(output output)

func Resp(ctx *gin.Context) ResultFunc {
	return func(result interface{}) func(output output) {
		return func(output output) {
			output(ctx, result)
		}
	}
}

func OK(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, v)
}

func ERR(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusBadRequest, v)
}

func SystemERR(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusInternalServerError, v)
}
