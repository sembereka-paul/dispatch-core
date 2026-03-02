### Development setup
#### Prerequisits
1. [Taskfile](https://taskfile.dev) Used as a task and script runner
2. Docker and Docker Compose are required for the development infrastructure
3. [Protoc](https://protobuf.dev/installation/) Protocol Buffers compiler for service-to-service codegen 
4. [goenv](https://github.com/go-nv/goenv?tab=readme-ov-file) Go version manager

#### Protobufs
This project uses gRPC for inter-service communication. The module `coop/proto` contains protobuf schemas and codegen.

From the repository root, run:
```
protoc --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative proto/event.proto
```
### Quick start
1. Clone the repository
```
git clone https://github.com/sembereka-paul/dispatch-core.git
```
2. Run `go mod tidy`. It's optional to run it in `api` or `pub-sub` directories only.
3. Run `task backend:up` or (just `task`) to run `pub-sub` and `api` in parallel.
3. `cd dashboard` to run the minimal ui. 
    - Run `npm ci` then `npm run dev` to run it in dev mode.

#### Running the services together
Run `task backend:up` or (just `task`) to run `pub-sub` and `api` in parallel

#### Pub-sub service - standalone
Run `task pub-sub:up` to run the `pub-sub` service.

The service depends on protobuf and [The Mastodon streaming api](https://docs.joinmastodon.org/methods/streaming/).

#### API service - standalone
Use `task api:up` to run the `api` service.

### Mastodon
#### User and auth token
1. Run `task dev:spin`. If it succeeds, run `docker ps` to find the container ID.
2. Create an admin user:
```
docker exec -it {container} bin/tootctl accounts create admin --email=admin@dev.local --confirmed --approved
```

The command returns the user password. Use the credentials to create an `access_token` and place it in .env.development as `MASTODON_ACCESS_TOKEN`. You can create additional users as needed.

#### Setup
Mastodon requires TLS and a local domain alias. Add a 127.0.0.1 entry in `/etc/hosts` and use that hostname as `MASTODON_BASE_URL` in .env.development and as `LOCAL_DOMAIN` in `.env.mastodon`. Without this mastodon will probably present some issues. Caddy is used as the reverse proxy to handle TLS certificates.

### Architecture

#### Project goal
The project simulates social media marketing campaigns by using Mastodon hashtags as campaign names. It provides infrastructure and services to collect and display posts for tracked hashtags. A minimal UI shows incoming posts; advanced analytics were not implemented.

#### Service design
Two services stream data between them:

1. API service: manages user subscriptions and streams to clients
2. Pub-sub service: subscribes to Mastodon’s streaming API, parses events, and forwards them to the API 

````
+-----------+     sse    +---------------+    grpc       +------------------+      sse      +----------+
| end-users | <----------|  API service  | <-----------> |  pub-sub service | <------------ | mastodon |
+-----------+            +---------------+               +------------------+               +----------+
                      
````
