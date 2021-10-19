package main

import (
	"github.com/PlatONE_Network/PlatONE-SDK-Go/ws"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	//ws.Test()
	gin.SetMode(gin.DebugMode)
	gracesRouter := ws.InitRouter()
	err := gracesRouter.Run("127.0.0.1:8888")
	if err != nil {
		logrus.Errorf("test start err: %v", err)
		return
	}

}

func test() {
	//ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	//c, _ := client.NewClient(ctx, "http://127.0.0.1:6791", "0", "./keystore")
	//client.RpcSend(c.RpcClient)

}
