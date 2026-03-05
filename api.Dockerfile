FROM golang:1.25-alpine as build
RUN apk add --no-cache git ca-certificates
WORKDIR /app

COPY proto ./proto
COPY api ./api
COPY pub-sub ./pub-sub
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./go.work ./go.work

RUN go work sync
RUN cd proto && go mod download
RUN cd api && go mod download
RUN go build -v -o ./bin/api ./api

FROM alpine:3.23.3 as runtime

ENV GO_ENV=production
ENV CGO_ENABLED=0
COPY --from=build /app/bin/api /app/bin/api
WORKDIR /app
EXPOSE 8080

ENTRYPOINT [ "./bin/api" ]
