# goAgent
## 生成 protoc 工具安装
* 下载地址：https://github.com/protocolbuffers/protobuf/releases
* windows 配置protoc 到path 目录
* go get -u github.com/golang/protobuf/protoc-gen-go 
* protoc -I inter/ inter/agent.proto --go_out=plugins=grpc:inter

