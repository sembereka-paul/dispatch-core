### Development setup
#### Prerequisits
1. [Taskfile](https://taskfile.dev) Uses as task and script runner
2. Docker dev infrastructure depends on docker and docker compose
3. [Protoc](https://protobuf.dev/installation/) for service-to-service communication
4. [goenv](https://github.com/go-nv/goenv?tab=readme-ov-file) gform olang version management

#### Pub-sub service
The service depens on protobuf. The module `coop/proto` manages protobus schemas and codegen.

Run the command below from root to for protobufs codegen.
```
protoc --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative coop/event.proto
```

Then run `go run coop/pub-sub` to run the `pub-sub` service.
