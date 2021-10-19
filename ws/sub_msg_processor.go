package ws

import (
	"errors"
	"fmt"
	"reflect"
)

func NewSubMsgProcessor() *SubMsgProcessor {
	return &SubMsgProcessor{}
}

// SubMsgProcessor 订阅消息处理
type SubMsgProcessor struct{}

func (s *SubMsgProcessor) Process(ctx *MsgProcessorContext, msg interface{}) error {

	return nil
}

func (s *SubMsgProcessor) process(ctx *MsgProcessorContext, topicName string, chainID string, params map[string]interface{}) error {
	methodName := MethodCapitalized(topicName)
	reType := reflect.TypeOf(s)
	method, ok := reType.MethodByName(methodName)
	if !ok {
		return errors.New(fmt.Sprintf("no process method for topic[%v]", topicName))
	}
	methodParams := make([]reflect.Value, 4)
	// 第一个参数为方法的持有者
	methodParams[0] = reflect.ValueOf(s)
	methodParams[1] = reflect.ValueOf(ctx.client)
	methodParams[2] = reflect.ValueOf(chainID)
	methodParams[3] = reflect.ValueOf(params)
	resValues := method.Func.Call(methodParams)
	if len(resValues) > 0 {
		if err, ok := resValues[len(resValues)-1].Interface().(error); ok {
			return err
		}
	}
	return nil
}

// NewHeads newHeads topic 的订阅数据处理方法，方法名称需要与 topic 保持一致，但首字母需要是大写的
func (s *SubMsgProcessor) NewHeads() error {
	return nil
}

func (s *SubMsgProcessor) forwardBlock() error {
	return nil
}

func (s *SubMsgProcessor) forwardTX() error {
	return nil
}

func (s *SubMsgProcessor) forwardStats() error {
	return nil
}

// Forward 把从链上接收到的事件数据转发到前端
func (s *SubMsgProcessor) Forward() error {
	return nil
}

// 获取统计数据
func (s *SubMsgProcessor) getStat(chainID string) {
}

func (s *SubMsgProcessor) getNodes(chainID string) error {

	return nil
}
