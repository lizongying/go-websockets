# go-websockets

golang 实现websockets长连接，可以进行简单的数据转发
[go-websockets](https://github.com/voiladata/go-websockets)

## Build

```shell
make
```

or

```shell
go vet ./cmd/server
go build -ldflags "-s -w" -o  ./assets/server  ./cmd/server

go vet ./cmd/client
go build -ldflags "-s -w" -o  ./assets/client  ./cmd/client
```

## Dev

```shell
go mod init github.com/lizongying/go-websockets
```

```shell
go run cmd/server/*.go --server=:1234 --path=/echo
go run cmd/server/*.go
```

```shell
go run cmd/client/*.go --url=ws://127.0.0.1:1234/echo --origin=http://127.0.0.1/
go run cmd/client/*.go --from=client1
go run cmd/client/*.go --from=client2 --to=client1 --msg=hi
go run cmd/client/*.go --from=client2 --via=client1 --to=client2 --msg=hi
go run cmd/client/*.go --from=client2 --via=client1 --to=client2 --msg=hi --wait
```