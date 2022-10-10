package middleware

import (
	"fmt"
	"github.com/asaswangc/gowork/middleware/session"
	"github.com/asaswangc/gowork/response"
	"github.com/asaswangc/gowork/result"
	"github.com/asaswangc/gowork/variable"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"strings"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch ty := err.(type) {
				case *result.ConstErr:
					response.Resp(ctx)(response.NewJsonResult(ty.GetCode(), fmt.Sprintf("%s", ty.Error()), ty.ErrComment))(response.OK)
				default:
					response.Resp(ctx)(response.NewJsonResult(variable.ResponseInternalErrCode, fmt.Sprintf("%v", ty), nil))(response.SystemERR)
				}
			}
		}()
		ctx.Next()
	}
}

func Authenticate(runMode string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch {
		case runMode != variable.ReleaseMode:
			ctx.Set("user_id", 1001)
			ctx.Set("tenant_id", 10013)
			ctx.Next()
		case strings.Contains(ctx.FullPath(), "/service"):
			if userId := cast.ToInt(ctx.Query("user_id")); userId == 0 {
				response.Resp(ctx)(response.NewJsonResult(result.AuthFailedErr.ErrCode, "请在查询参数中提供正确的用户ID", nil))(response.ERR)
				ctx.Abort()
				return
			} else {
				ctx.Set("user_id", userId)
			}
			if tenantId := cast.ToInt(ctx.Query("tenant_id")); tenantId == 0 {
				response.Resp(ctx)(response.NewJsonResult(result.AuthFailedErr.ErrCode, "请在查询参数中提供正确的租户ID", nil))(response.ERR)
				ctx.Abort()
				return
			} else {
				ctx.Set("tenant_id", tenantId)
			}
			ctx.Next()
		default:
			se, err := session.GetSessionData(ctx)
			if err != nil {
				response.Resp(ctx)(response.NewJsonResult(result.AuthFailedErr.ErrCode, err.Error(), nil))(response.ERR)
				ctx.Abort()
				return
			}
			if userId, ok := se["userId"].(float64); !ok {
				response.Resp(ctx)(response.NewJsonResult(result.AuthFailedErr.ErrCode, "用户未登录或session已过期", nil))(response.ERR)
				ctx.Abort()
				return
			} else {
				ctx.Set("user_id", int(userId))
			}
			if tenantId, ok := se["tenantId"].(float64); !ok {
				response.Resp(ctx)(response.NewJsonResult(result.AuthFailedErr.ErrCode, "用户未登录或session已过期", nil))(response.ERR)
				ctx.Abort()
				return
			} else {
				ctx.Set("tenant_id", int(tenantId))
			}
			ctx.Next()
		}
	}
}
