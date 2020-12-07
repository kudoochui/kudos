module github.com/kudoochui/kudos

go 1.15

require (
	github.com/OwnLocal/goes v1.0.0
	github.com/beego/goyaml2 v0.0.0-20130207012346-5545475820dd
	github.com/beego/x2j v0.0.0-20131220205130-a0352aadc542
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/kudoochui/rpcx v0.0.0-20201126113720-ea80572f0b9b
	github.com/mitchellh/mapstructure v1.3.3
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/prometheus/client_golang v1.8.0
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18
	github.com/wendal/errors v0.0.0-20181209125328-7f31f4b264ec // indirect
	google.golang.org/grpc/examples v0.0.0-20201125005357-44e408dab41e // indirect
	gotest.tools v2.2.0+incompatible
)

replace github.com/kudoochui/rpcx v0.0.0-20201126113720-ea80572f0b9b => ../rpcx
