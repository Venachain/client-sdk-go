# Go SDK 下载与使用

[TOC]

Venachain Go SDK是面向Go开发者，提供的Venachain联盟链的Go开发工具包，提供了在应用层（Go 代码）访问区块链节点并获取服务的接口，比如部署合约、调用合约、查询链上数据等。

## SDK 下载

请首先下载SDK最新版本的发布包，[下载地址](https://git-c.i.wxblockchain.com/vena/src/client-sdk-go)。

将发布包到本地目录，如下所示：

```shell
# 下载
git clone https://git-c.i.wxblockchain.com/vena/src/client-sdk-go.git
```

## SDK 使用

可以使用go mod replace的方式使用sdk。编辑go.mod 文件，添加以下内容：

```go
// 编辑go.mod 文件添加sdk 包
require git-c.i.wxblockchain.com/vena/src/client-sdk-go v0.0.0-00010101000000-000000000000
replace git-c.i.wxblockchain.com/vena/src/client-sdk-go => git clone 后sdk所在的文件路径
```

```go
// 导入sdk包
import (
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/client"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/client/asyn"
)
```

