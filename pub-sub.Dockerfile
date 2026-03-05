FROM golang:1.25-alpine as build

WORKDIR /app

COPY proto ./proto
COPY api ./api
COPY pub-sub ./pub-sub
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./go.work ./go.work

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
