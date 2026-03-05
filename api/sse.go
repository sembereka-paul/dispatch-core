package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Content string `json:"content"`
	Time    string `json:"time"`
}

func SSEHandler(c *gin.Context) {
	id := c.Param("id")

	if len(id) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad params"})
		return
	}

	client := &Client{make(chan Message)}
	registerClient(client)
	defer unregisterClient(client)

	camp := CampaignHandler(id, client)
	defer camp.MaybeUntrack(client)

	subTag <- id

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
