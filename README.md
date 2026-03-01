### Development setup
#### Prerequisits
1. [Taskfile](https://taskfile.dev) Used as a task and script runner
2. Docker dev infrastructure depends on docker and docker compose
3. [Protoc](https://protobuf.dev/installation/) for service-to-service communication
4. [goenv](https://github.com/go-nv/goenv?tab=readme-ov-file) for golang version management

#### Protobufs
The project uses grpc to communicate between services.
The module `coop/proto` manages protobus schemas and codegen.

Run the command below from root to for protobufs codegen.
```
protoc --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative coop/event.proto
```

#### Running the services together
Run `task backend:up` or just `task` to run `pub-sub` and `api` in parallel

#### Pub-sub service
Run `task pub-sub:up` to run the `pub-sub` service.

The service depends on protobuf and [The Mastodon streaming api](https://docs.joinmastodon.org/methods/streaming/).

#### API service
Use `task api:up` to run the `api` service.
