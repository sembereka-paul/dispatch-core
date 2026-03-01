package main

import "sync"

type Client struct {
	Channel chan Message
}

var (
	Clients      = make(map[*Client]bool)
	ClientsMutex sync.Mutex
)

func registerClient(client *Client) {
	ClientsMutex.Lock()
	Clients[client] = true
	ClientsMutex.Unlock()
}

func unregisterClient(client *Client) {
	ClientsMutex.Lock()
	delete(Clients, client)
	close(client.Channel)
	ClientsMutex.Unlock()
}
