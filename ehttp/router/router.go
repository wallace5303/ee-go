package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wallace5303/ee-go/ehelper"
	"github.com/wallace5303/ee-go/elog"

	"github.com/gin-gonic/gin"
)

// 处理函数
type HandlerFunc func(*Ctx)

type Ctx struct {
	GinCtx *gin.Context   // gin context
	Err    any            // 请求错误
	Timed  int64          // 执行时间
	Args   map[string]any // 请求参数
}

var (
	// gin
	GinRouter *gin.Engine
)

func SetGinRouter(router *gin.Engine) {
	GinRouter = router
}

func Handle(httpMethod, path string, handler HandlerFunc) {
	GinRouter.Handle(httpMethod, path, func(gc *gin.Context) {
		ctx := &Ctx{GinCtx: gc}
		begin := time.Now()
		bindErr := ctx.GinCtx.BindJSON(&ctx.Args)
		if bindErr != nil {
			ctx.Args = nil
			//elog.CoreLogger.Errorf("bind args err: %+v", ctx.Args)
		}

		defer func() {
			ctx.Timed = time.Since(begin).Milliseconds()
			if err := recover(); err != nil {
				ctx.Err = err
			}

			record := fmt.Sprintf("[ee-go] http_method:%s, path:%s, req_params:%+v, exec_time:%dms", httpMethod, path, ctx.Args, ctx.Timed)
			elog.CoreLogger.Infof(record)
		}()
		handler(ctx)
	})
}

func (c *Ctx) JSON(data any) {
	c.GinCtx.JSON(http.StatusOK, data)
}

func (c *Ctx) JSONWithCode(code int, data any) {
	c.GinCtx.JSON(code, data)
}

func (c *Ctx) ArgJson() (arg map[string]any, ok bool) {
	result := ehelper.GetJson()
	arg = c.Args

	if len(arg) == 0 {
		result.Code = -1
		result.Msg = "parses request failed"
		return
	}

	ok = true
	return
}
