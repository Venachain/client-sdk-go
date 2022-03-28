package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/ws"
)

var (
	// DefaultWSSubscriber 默认的 websocket 订阅器
	DefaultWSSubscriber *ws.WsSubscriber
)

func InitWsSubscriber() {
	logrus.Debugf("DefaultWSSubscriber init [start]")
	DefaultWSSubscriber = ws.NewWSSubscriber("127.0.0.1", 26791, "platone")
	logrus.Debugf("DefaultWSSubscriber init [end]")
}

// 以下为websocket 测试
func main() {
	InitWsSubscriber()
	gin.SetMode(gin.DebugMode)
	gracesRouter := ws.InitRouter()
	err := gracesRouter.Run("127.0.0.1:8888")
	if err != nil {
		logrus.Errorf("test start err: %v", err)
		return
	}
}
