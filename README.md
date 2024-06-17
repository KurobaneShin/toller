# toller

```
docker compose up -d
```

## Installing protobuf compiler
```
sudo apt install -y protobuf-compiler
```

## Installing GRPC and Protobuffer plugins for Golang
1. Protobuffers
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@1.28
```
2. GRPC
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
3. install the package dependencies
```
go mod tidy
```
or
```
go get google.golang.org.protobuf
go get google.golang.org/grpc
```
