package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

var (
	// DefaultWSSubscriber 默认的 websocket 订阅器
	DefaultWSSubscriber *wsSubscriber
)

func init() {
	logrus.Debugf("DefaultWSSubscriber init [start]")
	DefaultWSSubscriber = newWSSubscriber()
	logrus.Debugf("DefaultWSSubscriber init [end]")
}

type wsSubscriber struct {
	wsManager *Manager
	ip        string
	port      int64
	//path      string
	group string
}

func newWSSubscriber() *wsSubscriber {
	return &wsSubscriber{
		wsManager: DefaultWebsocketManager,
		ip:        "127.0.0.1",
		port:      26791,
		group:     "platone",
	}
}

// SubHeadForChain 为指定的链订阅它所配置的所有 websocket topics
func (s *wsSubscriber) SubHeadForChain() error {
	// 1、获取 websocket 客户端连接
	client, err := s.getWSClientByChain()
	if err != nil {
		return err
	}

	// 2、提取当前链订阅的 topics 配置信息
	topics := getNewHeadTopic()
	logrus.Debug(topics)

	// 3 处理 topics 订阅
	for _, topic := range topics {
		topicName := topic.(map[string]interface{})["name"]
		topicParams := topic.(map[string]interface{})["params"]
		topics := []string{"newHeads"}
		err := s.wsSubTopicProcessor(client, topicName.(string), topicParams.(string), topics)
		if err != nil {
			logrus.Warningln(err)
		}
	}
	return nil
}

// SubHeadForChain 为指定的链订阅它所配置的所有 websocket topics
func (s *wsSubscriber) SubLogForChain(address, topic string) error {
	// 1、获取 websocket 客户端连接
	client, err := s.getWSClientByChain()
	if err != nil {
		return err
	}

	// 2、提取当前链订阅的 topics 配置信息
	topics := getLogTopic(address, topic)
	logrus.Debug(topics)

	// 3 处理 topics 订阅
	for _, topic := range topics {
		topicName := topic.(map[string]interface{})["name"]
		topicParams := topic.(map[string]interface{})["params"]
		topics := []string{"log"}
		err := s.wsSubTopicProcessor(client, topicName.(string), topicParams.(string), topics)
		if err != nil {
			logrus.Warningln(err)
		}
	}
	return nil
}

// topic 订阅处理器
func (s *wsSubscriber) wsSubTopicProcessor(client *Client, topic string, params string, topics []string) error {
	// 获取到该链所配置订阅的所有 topic

	exist := false
	for _, v := range topics {
		if topic == v {
			exist = true
			break
		}
	}
	if !exist {
		return errors.New(fmt.Sprintf("unknown websocket subscription topic: %v", topic))
	}
	methodName := MethodCapitalized(topic)
	// 通过反射进行方法调用
	reType := reflect.TypeOf(s)
	method, ok := reType.MethodByName(methodName)
	if !ok {
		return errors.New(fmt.Sprintf("no process method for topic[%v]", topic))
	}
	methodParams := make([]reflect.Value, 4)
	// 第一个参数为方法的持有者
	methodParams[0] = reflect.ValueOf(s)
	//methodParams[1] = reflect.ValueOf(chain)
	methodParams[1] = reflect.ValueOf(client)
	methodParams[2] = reflect.ValueOf(topic)
	methodParams[3] = reflect.ValueOf(params)
	resValues := method.Func.Call(methodParams)
	if len(resValues) > 0 {
		if err, ok := resValues[len(resValues)-1].Interface().(error); ok {
			return err
		}
	}
	return nil
}

// NewHeads newHeads 事件的订阅处理
func (s *wsSubscriber) NewHeads(client *Client, topic string, params string) error {
	// 1、处理参数信息
	paramsStr, err := s.wsSubParamsProcess(topic, params)
	if err != nil {
		return err
	}

	// 2、给服务端发送订阅消息对指定 topic 进行订阅
	dto := WSMessageDTO{
		ID:      client.Id,
		Group:   client.Group,
		Message: paramsStr,
	}
	s.wsManager.Send(dto.ID, dto.Group, []byte(dto.Message))
	logrus.Infof("subscribe topic[newHead] from websocket for chain[%v] [success]", "platone")
	return nil
}

// Log log 事件的订阅处理
func (s *wsSubscriber) Log(client *Client, topic string, params string) error {
	// 1、处理参数信息
	paramsStr, err := s.wsSubParamsProcess(topic, params)
	if err != nil {
		return err
	}

	// 2、给服务端发送订阅消息对指定 topic 进行订阅
	dto := WSMessageDTO{
		ID:      client.Id,
		Group:   client.Group,
		Message: paramsStr,
	}
	s.wsManager.Send(dto.ID, dto.Group, []byte(dto.Message))
	logrus.Infof("subscribe topic[newHead] from websocket for chain[%v] [success]", "platone")
	return nil
}

// 处理 websocket 事件订阅的参数
func (s *wsSubscriber) wsSubParamsProcess(topic string, params string) (string, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(params), &data)
	if err != nil {
		logrus.Errorf("websocket params unmarshal error: %v", err)
		return "", err
	}
	msgType := data["id"]
	// type topic id
	data["id"] = fmt.Sprintf("%v %v", msgType, topic)

	p, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("websocket params marshal error: %v", err)
		return "", err
	}
	params = string(p)
	return params, nil
}

func (s *wsSubscriber) getWSClientByChain() (*Client, error) {
	ip := s.ip
	port := s.port
	//path := s.path
	group := s.group
	client, err := s.wsManager.Dial(ip, port, group)
	url := fmt.Sprintf("ws://%s:%v", ip, port)
	if err != nil {
		msg := fmt.Sprintf("chain[%s][%s:%v] websocket dial [%s] error: %v",
			group, ip, port, url, err)
		return nil, errors.New(msg)
	}
	return client, nil
}

func getNewHeadTopic() map[string]interface{} {
	topic := make(map[string]interface{})
	tmp := make(map[string]interface{}, 2)
	tmp["name"] = "newHeads"
	tmp["params"] = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_subscribe\", \"params\": [\"newHeads\"],\"id\":\"subscription\"}"
	topic["new_heads"] = tmp
	return topic
}

func getLogTopic(address, topics string) map[string]interface{} {
	topic := make(map[string]interface{})
	tmp := make(map[string]interface{}, 2)
	tmp["name"] = "log"
	tmp["params"] = "{\"id\": \"subscription\", \"method\": \"eth_subscribe\", \"params\": [\"logs\", {\"address\":\"" + address + "\", \"topics\": [\"" + topics + "\"]}]}"
	topic["new_heads"] = tmp
	return topic
}
