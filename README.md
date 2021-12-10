# PlatONE-SDK-GO
PlatONE Go SDK是面向go开发者，提供的PlatONE联盟链的go开发工具包，提供了在应用层访问区块链节点并获取服务的接口，比如部署合约、调用合约、查询链上数据等。


## SDK 下载

**Golang Get 添加私有仓库, 通过`ssh`公钥访问私有仓库**

打开 ~/.gitconfig 文件 添加以下内容
```text
[url "ssh://git@git-c.i.wxblockchain.com/"]
    insteadOf = https://git-c.i.wxblockchain.com/
```

添加go私有仓库
```bash
export GOPRIVATE=git-c.i.wxblockchain.com
go get git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go
```

**对于gitlab多级子group的情况，需要设置如下**

```bash
# 新建~/.netrc文件，配置git-c.i.wxblockchain.com登录信息
# username: OA账号
# git_password: access_token 或者 OA密码
machine git-c.i.wxblockchain.com login username password git_password
```

修改依赖包获取源
```bash
export GOPROXY=https://goproxy.cn,direct
```

