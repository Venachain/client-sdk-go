package ws

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	myRouter *gin.Engine
)

func init() {
	myRouter = gin.Default()
}

func InitRouter() *gin.Engine {
	myRouter.Use(PanicHandler())
	myRouter.Use(Cors())

	customRouter()

	return myRouter
}

// PanicHandler gin 统一 panic 处理器
func PanicHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("unknown panic：%+v", err)
				ctx.Abort()
				return
			}
		}()
		ctx.Next()
	}
}

// Cors 处理跨域资源共享问题
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "127.0.0.1")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func customRouter() {
	api := myRouter.Group("/api")
	{
		// websocket
		wsGroup := api.Group("/ws")
		{
			if gin.Mode() == gin.DebugMode {
				wsGroup.StaticFile("/ws_sub_test.html", "./ws_sub_test.html")
			}

			wsGroup.GET("/", DefaultWebsocketManager.WsClient)
		}
	}
}

func Test() {
	DefaultWSSubscriber.SubTopicsForChain()
	time.Sleep(100 * time.Second)
}
