package client

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/sirupsen/logrus"
)

type WebSocket struct {
	IP     string
	WsPort uint64
}

// Dial 作为 websocket 客户端拨号去连接其他 websocket 服务端
func (ws WebSocket) subcribe() (interface{}, error) {

	url := fmt.Sprintf("ws://%s:%v", ws.IP, ws.WsPort)
	conn, resp, err := websocket.DefaultDialer.DialContext(context.Background(), url, nil)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("websocket dial success, response: %+v", resp)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		msg := string(message)
		fmt.Println(msg)
		//logrus.Debugf("client [%s] receive message: %s", c.Id, msg)
		//go c.readMessageProcessor(msg)
	}

	return nil, nil
}
