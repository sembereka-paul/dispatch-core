package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	pb "coop/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Message struct {
	Content string `json:"content"`
	Time    string `json:"time"`
}

type Client struct {
	Channel chan Message
}

type Campaign struct {
	Id        string           // tag or campaign id
	Listeners map[*Client]bool // campaign subscribers
}

var (
	clients            = make(map[*Client]bool)
	clientsMutex       sync.Mutex
	campaigns          = make(map[string]*Campaign)
	campaignsMutex     sync.Mutex
	subscriptions      = make(map[string]int)
	subscriptionsMutex sync.Mutex
)

// NotificationChannel holds channels for clients to receive updates
var NotificationChannel = make(chan Message)

func SSEHandler(c *gin.Context) {
	id := c.Param("id")

	if len(id) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad params"})
		return
	}

	client := &Client{make(chan Message)}
	registerClient(client)
	defer unregisterClient(client)

	camp := getCampaign(id)
	if camp == nil {
		camp = &Campaign{id, map[*Client]bool{}}
	}
	trackCampaign(camp, client)
	defer camp.maybeUntrack(client)

	log.Println("attempt sub", id)
	subTag <- id
	log.Println("added sub", id)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	for {
		select {
		case msg := <-client.Channel:
			fmt.Fprintf(c.Writer, "data: %s\n\n", msg.Content)
			c.Writer.Flush()
		case <-c.Writer.CloseNotify():
			return
		}
	}
}

func (campaign *Campaign) AddListener(client *Client) {
	found := false
	for entry := range campaign.Listeners {
		if reflect.DeepEqual(entry, client) {
			found = true
			break
		}
	}
	if !found {
		campaign.Listeners[client] = true
	}
}

func addSubscription(name string) {
	log.Println("Add subsription req:", name)
	subscriptionsMutex.Lock()
	if count, ok := subscriptions[name]; ok {
		log.Println("Incrementing subsription:", name)
		subscriptions[name] = count + 1
	} else {
		log.Println("Adding subsription:", name)
		subscriptions[name] = 1
	}
	subscriptionsMutex.Unlock()
}

func getSubscription(name string) string {
	sub := ""
	subscriptionsMutex.Lock()
	if count, ok := subscriptions[name]; ok && count > 0 {
		sub = name
	}
	subscriptionsMutex.Unlock()
	return sub
}

func maybeRemoveSubscription(name string) {
	subscriptionsMutex.Lock()
	if count, ok := subscriptions[name]; ok && (count-1) == 0 {
		log.Println("Removing subsription:", name)
		delete(subscriptions, name)
	}
	subscriptionsMutex.Unlock()
}
func trackCampaign(campaign *Campaign, client *Client) {
	campaignsMutex.Lock()
	campaign.AddListener(client)
	campaigns[campaign.Id] = campaign
	campaignsMutex.Unlock()
	log.Println("Tracking campaign clients:", campaign.Id)

}

func (campaign *Campaign) maybeUntrack(client *Client) {
	campaignsMutex.Lock()
	if len(campaign.Listeners) == 1 && campaign.Listeners[client] {

		log.Println("Untracking campaign:", campaign.Id)
		delete(campaigns, campaign.Id)
	} else {
		log.Println("Untracking campaign clients:", campaign.Id)
		delete(campaign.Listeners, client)
	}
	campaignsMutex.Unlock()
}

func getCampaign(id string) *Campaign {
	var camp *Campaign
	campaignsMutex.Lock()
	camp = campaigns[id]
	campaignsMutex.Unlock()
	return camp
}

func registerClient(client *Client) {
	clientsMutex.Lock()
	clients[client] = true
	clientsMutex.Unlock()
}

func unregisterClient(client *Client) {
	clientsMutex.Lock()
	delete(clients, client)
	close(client.Channel)
	clientsMutex.Unlock()
}

func (campaign *Campaign) broadcast(msg Message) {
	clientsMutex.Lock()
	for client := range campaign.Listeners {
		client.Channel <- msg
	}
	clientsMutex.Unlock()
}

type server struct {
	pb.UnimplementedEventServer
}

type EventMessage struct {
	event string `json:"event"`
	data  string `data:"data"`
	tag   string `json:"tag"`
}

var see = make(chan EventMessage, 1)
var subTag = make(chan string)

// subscribe
func sub(c pb.EventClient, tag string) error {
	defer maybeRemoveSubscription(tag)

	addSubscription(tag)

	fmt.Println("sub", tag)
	ctx := context.Background()
	defer ctx.Done()

	stream, err := c.Sub(ctx, &pb.SubscribeRequest{
		Tag: tag,
	})
	if err != nil {
		log.Println("Err with stream", err)
		return err
	}

	for {
		event, err := stream.Recv()

		if err != nil {
			continue
		}
		see <- EventMessage{
			event: event.Event,
			data:  event.Data,
			tag:   event.Tag,
		}
	}

}

func main() {

	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEventClient(conn)

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-see:
				log.Println("received:", msg)
				camp := getCampaign(msg.tag)
				if camp != nil {
					camp.broadcast(Message{
						Time:    time.Now().Format(time.RFC3339),
						Content: msg.data,
					})
				}
			case <-quit:
				return
			}
		}
	}()

	go func() {
		for subReq := range subTag {
			log.Println("Attempt go sub:", subReq)
			go func(tag string) {
				subscription := getSubscription(subReq)
				if subscription == "" {
					if err := sub(c, tag); err != nil {
						log.Println("sub error:", err)
					}
				}
			}(subReq)
		}
	}()

	r := gin.Default()
	r.GET("/notifications/:id", SSEHandler)
	r.Run(":8080")
}
