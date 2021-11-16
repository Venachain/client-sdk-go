package main

import (
	"github.com/PlatONE_Network/PlatONE-SDK-Go/ws"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 以下为websocket 测试
func main() {
	gin.SetMode(gin.DebugMode)
	gracesRouter := ws.InitRouter()
	err := gracesRouter.Run("127.0.0.1:8888")
	if err != nil {
		logrus.Errorf("test start err: %v", err)
		return
	}
}
