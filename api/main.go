package main

import (
	"log"
	"os"
	"time"

	pb "coop/proto"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedEventServer
}

type EventMessage struct {
	event string `json:"event"`
	data  string `data:"data"`
	tag   string `json:"tag"`
}

var eventMessage = make(chan EventMessage, 1)
var subTag = make(chan string)

var (
	ENV = os.Getenv("ENV")
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEventClient(conn)

	go subscriptionsManager(c)
	go messageDispatcher()

	r := gin.Default()

	if ENV == "development" {
		r.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: false,
		}))
	}
	r.GET("/notifications/:id", SSEHandler)
	r.Run(":8080")
}

func subscriptionsManager(connection pb.EventClient) {
	for subReq := range subTag {
		log.Println("Attempt sub:", subReq)
		go func(tag string) {
			subscription := getSubscription(subReq)
			if subscription == "" {
				if err := sub(connection, tag); err != nil {
					log.Println("sub error:", err)
				}
			}
		}(subReq)
	}
}

func messageDispatcher() {
	for msg := range eventMessage {
		log.Println("received:", msg)
		camp := GetCampaign(msg.tag)
		if camp != nil {
			camp.Broadcast(Message{
				Time:    time.Now().Format(time.RFC3339),
				Content: msg.data,
			})
		}
	}
}
