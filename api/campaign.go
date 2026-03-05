package main

import (
	"log"
	"reflect"
	"sync"
)

type Campaign struct {
	Id        string           // tag or campaign id
	Listeners map[*Client]bool // campaign subscribers
}

var (
	Campaigns      = make(map[string]*Campaign)
	CampaignsMutex sync.RWMutex
)

func CampaignHandler(tag string, client *Client) *Campaign {
	CampaignsMutex.Lock()
	defer CampaignsMutex.Unlock()

	camp := Campaigns[tag]
	if camp == nil {
		camp = &Campaign{tag, map[*Client]bool{}}
	}

	camp.AddListener(client)
	Campaigns[camp.Id] = camp
	log.Println("Tracking campaign clients:", camp.Id)
	return camp
}

func (campaign *Campaign) MaybeUntrack(client *Client) {
	CampaignsMutex.Lock()
	defer CampaignsMutex.Unlock()

	if len(campaign.Listeners) == 1 && campaign.Listeners[client] {

		log.Println("Untracking campaign:", campaign.Id)
		delete(Campaigns, campaign.Id)
	} else {
		log.Println("Untracking campaign clients:", campaign.Id)
		delete(campaign.Listeners, client)
	}
}

func GetCampaign(id string) *Campaign {
	CampaignsMutex.RLock()
	defer CampaignsMutex.RUnlock()

	return Campaigns[id]
}
func (campaign *Campaign) Broadcast(msg Message) {
	ClientsMutex.RLock()
	defer ClientsMutex.RUnlock()

	for client := range campaign.Listeners {
		client.Channel <- msg
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
