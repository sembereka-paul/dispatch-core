package main

import "sync"

type Client struct {
	Channel chan Message
}

var (
	Clients      = make(map[*Client]bool)
	ClientsMutex sync.RWMutex
)

func registerClient(client *Client) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()
	Clients[client] = true
}

func unregisterClient(client *Client) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	delete(Clients, client)
	close(client.Channel)
}
