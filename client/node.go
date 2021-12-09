package client

import (
	"context"
	"encoding/json"
	"errors"

	common_sdk "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/common"
	precompile "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled/syscontracts"
)

type NodeClient struct {
	ContractClient
	NodeName string
}

// 必传参数为<publicKey>: 节点公钥，用于节点间安全通信。节点的公私钥对可由ethkey工具产生。
//<externalIP>: 节点外网IP
//<internalIP>: 节点内网IP
func (nodeClient NodeClient) NodeAdd(ctx context.Context, txparam common_sdk.TxParams, requestNodeInfo syscontracts.NodeInfo) (string, error) {
	funcName := "add"
	nodeInfo, err := setNodeInfoDefault(nodeClient, requestNodeInfo)
	if err != nil {
		return "", err
	}
	bytes, _ := json.Marshal(nodeInfo)
	strJson := string(bytes)
	funcParams := []string{strJson}

	result, err := nodeClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.NodeManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 将节点从节点列表中删除。在下一次peers更新后，被删除的节点会被PlatONE网络中的其他节点断开连接。
func (nodeClient NodeClient) NodeDelete(ctx context.Context, txparam common_sdk.TxParams) (string, error) {
	funcName := "update"
	var str = "{\"status\":2}"

	if err := common_sdk.ParamValid(nodeClient.NodeName, "name"); err != nil {
		return "", err
	}
	funcParams := common_sdk.CombineFuncParams(nodeClient.NodeName, str)

	result, err := nodeClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.NodeManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (nodeClient NodeClient) NodeUpdate(ctx context.Context, txparam common_sdk.TxParams, request syscontracts.NodeUpdateInfo) (string, error) {
	funcName := "update"
	if err := common_sdk.ParamValid(nodeClient.NodeName, "name"); err != nil {
		return "", err
	}
	bytes, _ := json.Marshal(request)
	strJson := string(bytes)
	funcParams := []string{nodeClient.NodeName, strJson}
	result, err := nodeClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.NodeManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 如果传入为nil，则查询所有
func (nodeClient NodeClient) NodeQuery(ctx context.Context, txparam common_sdk.TxParams, request *syscontracts.NodeQueryInfo) (string, error) {
	if request == nil {
		funcName := "getAllNodes"
		result, err := nodeClient.contractCallWrap(ctx, txparam, nil, funcName, precompile.NodeManagementAddress)
		if err != nil {
			return "", err
		}
		res := result.([]interface{})
		return res[0].(string), nil
	}
	funcName := "getNodes"
	if err := common_sdk.ParamValid(nodeClient.NodeName, "name"); err != nil {
		return "", err
	}
	bytes, _ := json.Marshal(request)
	funcParams := []string{string(bytes)}

	result, err := nodeClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.NodeManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 通过查询键对节点信息进行查询，对匹配成功的数据对象进行统计，返回统计值，如果不需要，则传入为nil
func (nodeClient NodeClient) NodeStat(ctx context.Context, txparam common_sdk.TxParams, request *syscontracts.NodeStatInfo) (int32, error) {
	funcName := "nodesNum"
	var funcParams []string
	m := make(map[string]interface{})
	if err := common_sdk.ParamValid(nodeClient.NodeName, "name"); err != nil {
		return 0, err
	}
	if request == nil {
		return 0, errors.New("parameter is incorrect\n")
	}
	// Status 为1有效状态或者2无效状态，如果status 为0，则表示传入的参数为构造参数时设置的默认值，以下为处理默认值的逻辑
	if request.Status == 0 {
		m["type"] = request.Type
		bytes, _ := json.Marshal(m)
		funcParams = []string{string(bytes)}
	} else {
		bytes, _ := json.Marshal(request)
		funcParams = []string{string(bytes)}
	}
	result, err := nodeClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.NodeManagementAddress)
	//cxh
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(int32), nil
}

func setNodeInfoDefault(nodeClient NodeClient, requestNodeInfo syscontracts.NodeInfo) (*syscontracts.NodeInfo, error) {
	if requestNodeInfo.ExternalIP == "" || requestNodeInfo.InternalIP == "" || requestNodeInfo.PublicKey == "" {
		return nil, errors.New("insufficient parameters")
	}
	requestNodeInfo.Name = nodeClient.NodeName
	requestNodeInfo.Status = 1
	if requestNodeInfo.RpcPort == 0 {
		requestNodeInfo.RpcPort = 6791
	}
	if requestNodeInfo.P2pPort == 0 {
		requestNodeInfo.P2pPort = 16791
	}
	return &requestNodeInfo, nil
}
