package ws

import (
	"encoding/json"
	"errors"
	"fmt"
)

func NewSubMsgProcessor() *SubMsgProcessor {
	return &SubMsgProcessor{}
}

// SubMsgProcessor 订阅消息处理
type SubMsgProcessor struct{}

func (s *SubMsgProcessor) Process(ctx *MsgProcessorContext, msg interface{}) error {
	data, ok := msg.(map[string]interface{})
	if !ok {
		errStr := fmt.Sprintf("can't process unknown subMsg:\n %+v", msg)
		return errors.New(errStr)
	}

	jsonMsg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	DefaultWebsocketManager.SendGroup(ClientGroup, jsonMsg)

	return nil
}
