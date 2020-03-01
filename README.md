# kudos
kudos是一款基于微服务架构的简洁，高性能，易扩展，易部署的分布式游戏服务框架。基于rpcx的rpc，支持pomelo通信协议，轻松应用于游戏开发。


## 特点
- **简单**：容易上手，游戏开发需要基本组件和服务都已集成，直接调用。对于熟悉pomelo的特别友好。
- **组件化**：功能分为一个个组件，按需要加载。
- **分布式**：可以分成多个节点分布式部署，也可以打包一起作为一个进程部署。
- **微服务架构，支持服务发现**：consul，etcd，zookeeper等主流注册中心。
- **基于rpcx的rpc**：rpcx是一款高性能的rpc框架。其性能远远高于 Dubbo、Motan、Thrift等框架，是gRPC性能的两倍。支持服务治理。更多功能请参考：[http://rpcx.io](http://rpcx.io/)
- **跨语言**：除go外，还可以访问其它语言实现的节点服务。得益于rpcx。
- **支持pomelo通信协议**：该协议广泛用于各种游戏开发中，支持多端，多种语言版本。
- **易部署**：各服务器独立，无依赖，可以单独启动。

## 安装

`go get -u -v github.com/kudoochui/kudos`

## 开发脚手架(示例)
[kudosServer](https://github.com/kudoochui/kudosServer)

## 游戏架构参考
[游戏微服务架构设计：MMORPG](https://www.toutiao.com/i6798800455955644935/)

[游戏微服务架构设计：挂机类游戏](https://www.toutiao.com/i6798814918574342660/)

[游戏微服务架构设计：棋牌游戏](https://www.toutiao.com/i6798815085935460876/)

[游戏微服务架构设计：io游戏](https://www.toutiao.com/i6798815271386612231/)

## Roadmap
功能持续开发中

## 交流
[wiki](https://github.com/kudoochui/kudos/wiki)

QQ交流群：77584553

关注头条号：丁玲隆咚呛
分享更多内容

## 证书
MIT License