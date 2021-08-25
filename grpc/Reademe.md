官方连接：https://www.grpc.io/docs/languages/go/

学习基本的gRPC, 流式RPC

proto编译环境安装:
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
export PATH="$PATH:$(go env GOPATH)/bin"

proto文件编译命令:
route_guide: protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     route_guide/route_guide.proto
