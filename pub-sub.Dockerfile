FROM golang:1.25-alpine as build

# DO protobufs
RUN apk update && apk add --no-cache make protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1 \
    && export PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /app

COPY proto ./proto
COPY api ./api
COPY pub-sub ./pub-sub
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./go.work ./go.work

RUN protoc --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative proto/event.proto

RUN go work sync
RUN cd proto && go mod download
RUN cd pub-sub && go mod download
RUN go build -v -o ./bin/pub-sub ./pub-sub

FROM alpine:3.23.3 as runtime

ENV GO_ENV=production
ENV CGO_ENABLED=0
COPY --from=build /app/bin/pub-sub /app/bin/pub-sub
WORKDIR /app
EXPOSE 50052

ENTRYPOINT [ "./bin/pub-sub" ]
