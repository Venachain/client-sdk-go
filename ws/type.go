package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Manager 所有 websocket 信息
type Manager struct {
	Group                   map[string]map[string]*Client
	groupCount, clientCount uint
	Lock                    sync.Mutex
	Register, UnRegister    chan *Client
	Message                 chan *MessageData
	GroupMessage            chan *GroupMessageData
	BroadCastMessage        chan *BroadCastMessageData
}

// Client 单个 websocket 信息
type Client struct {
	Id, Group  string
	LocalAddr  string
	RemoteAddr string
	Path       string
	Socket     *websocket.Conn
	IsAlive    bool
	IsDial     bool
	RetryCnt   int64
	Message    chan []byte
}

// MessageData 单个客户端发送数据信息
type MessageData struct {
	Id, Group string
	Message   []byte
}

// GroupMessageData 组客户端广播数据信息
type GroupMessageData struct {
	Group   string
	Message []byte
}

// BroadCastMessageData 所有客户端广播数据信息
type BroadCastMessageData struct {
	Message []byte
}

// WSMessageDTO websocket 单个客户端发送数据信息
type WSMessageDTO struct {
	// websocket 客户端连接 ID
	ID string `json:"id" binding:"required"`
	// websocket 客户端连接所在的分组
	Group string `json:"group" binding:"required"`
	// 要发送的消息内容
	Message string `json:"message" binding:"required"`
}

// WSSubMsgDTO 与前端 ws 连接交互所需的消息 dto
type WSSubMsgDTO struct {
	// 消息id
	ID string `json:"id"`
	// 消息类型
	Type string `json:"type"`
	// 消息内容
	Content interface{} `json:"content"`
}

type Subscription interface {
	SubHeadForChain() error
	SubLogForChain(address, topic string) error
}
