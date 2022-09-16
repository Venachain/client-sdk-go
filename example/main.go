package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Venachain/client-sdk-go/ws"
)

// 以下为websocket 测试
func main() {
	ws.InitWsSubscriber("127.0.0.1", 26791, "venachain")
	gin.SetMode(gin.DebugMode)
	gracesRouter := ws.InitRouter()
	err := gracesRouter.Run("127.0.0.1:8888")
	if err != nil {
		logrus.Errorf("test start err: %v", err)
		return
	}
}
