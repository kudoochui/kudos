# Kudos
Kudos is a simple, high-performance, easy to expand and easy to deploy distributed game service framework
based on microservice architecture, It is based on RPC of rpcx, supports pomelo communication protocol and
can be easily applied to game development.

[中文文档](Chinese)

## Features
-**Easy to use**: Game development requires that basic components and services have been integrated and called directly. Especially friendly to those who are familiar with pomelo.

-**Componentization**: The functions are divided into components and loaded as required.

-**Distributed**: It can be deployed in multiple nodes or packaged together as a process.

-**Microservice architecture, supporting service discovery**: Mainstream registries such as consult, etcd, zookeeper, etc.

-**RPC based on rpcx**: rpcx is a high-performance RPC framework. Its performance is much higher than Dubbo, Motan, thrift and other frameworks, which is twice the performance of grpc. Support service governance. For more functions, please refer to:[ http://rpcx.io ] http://rpcx.io/ )

-**Cross language**: In addition to go, you can also access node services implemented in other languages. Thanks to rpcx.

-**Support pomelo communication protocol**: The protocol is widely used in various game development, supporting multi terminal and multi language versions.

-**Easy to deploy**: Each server is independent and independent and can be started independently.

## Installation
`go get -u -v github.com/kudoochui/kudos`

## Getting started(开发脚手架)
[kudosServer](https://github.com/kudoochui/kudosServer)

## 游戏架构参考
[游戏微服务架构设计：MMORPG](https://www.toutiao.com/i6798800455955644935/)

[游戏微服务架构设计：挂机类游戏](https://www.toutiao.com/i6798814918574342660/)

[游戏微服务架构设计：棋牌游戏](https://www.toutiao.com/i6798815085935460876/)

[游戏微服务架构设计：io游戏](https://www.toutiao.com/i6798815271386612231/)

## Roadmap
- Add more connector
- Actor support

## Community
[wiki](https://github.com/kudoochui/kudos/wiki)

QQ交流群：77584553

关注头条号：丁玲隆咚呛
分享更多内容

## 证书
MIT License