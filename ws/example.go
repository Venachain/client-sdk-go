package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

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
				wsGroup.StaticFile("/ws_sub_test.html", "./ws/ws_sub_test.html")
			}
			wsGroup.GET("/log/:group", DefaultWebsocketManager.WsClientForLog)
			wsGroup.GET("/head/:group", DefaultWebsocketManager.WsClientForNewHeads)
		}
	}
}

// WsClient gin 处理 websocket handler
func (manager *Manager) WsClientForLog(ctx *gin.Context) {
	group := ctx.Param("group")
	ClientGroup = group
	// 创建 websocket 升级器
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}
	// 将 http 升级为 websocket
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Errorf("websocket connect error: %s", group)
		return
	}

	client := &Client{
		Id:         uuid.NewV4().String(),
		Group:      group,
		LocalAddr:  conn.LocalAddr().String(),
		RemoteAddr: conn.RemoteAddr().String(),
		Path:       ctx.Request.URL.String(),
		Socket:     conn,
		IsAlive:    true,
		IsDial:     false,
		RetryCnt:   0,
		Message:    make(chan []byte, BuffSize),
	}
	manager.RegisterClient(client)
	address := "0x1000000000000000000000000000000000000005"
	topic := "0x8cd284134f0437457b5542cb3a7da283d0c38208c497c5b4b005df47719f98a1"
	if err = DefaultWSSubscriber.SubLogForChain(address, topic); err != nil {
		fmt.Errorf("subTopicsForChain is error")
	}
	go client.Read()
	go client.Write()
}

// WsClient gin 处理 websocket handler
func (manager *Manager) WsClientForNewHeads(ctx *gin.Context) {
	group := ctx.Param("group")
	ClientGroup = group
	// 创建 websocket 升级器
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}
	// 将 http 升级为 websocket
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Errorf("websocket connect error: %s", group)
		return
	}

	client := &Client{
		Id:         uuid.NewV4().String(),
		Group:      group,
		LocalAddr:  conn.LocalAddr().String(),
		RemoteAddr: conn.RemoteAddr().String(),
		Path:       ctx.Request.URL.String(),
		Socket:     conn,
		IsAlive:    true,
		IsDial:     false,
		RetryCnt:   0,
		Message:    make(chan []byte, BuffSize),
	}
	manager.RegisterClient(client)
	if err = DefaultWSSubscriber.SubHeadForChain(); err != nil {
		fmt.Errorf("subTopicsForChain is error")
	}
	go client.Read()
	go client.Write()
}
