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
	CampaignsMutex sync.Mutex
)

func TrackCampaign(campaign *Campaign, client *Client) {
	CampaignsMutex.Lock()
	campaign.AddListener(client)
	Campaigns[campaign.Id] = campaign
	CampaignsMutex.Unlock()
	log.Println("Tracking campaign clients:", campaign.Id)

}

func (campaign *Campaign) MaybeUntrack(client *Client) {
	CampaignsMutex.Lock()
	if len(campaign.Listeners) == 1 && campaign.Listeners[client] {

		log.Println("Untracking campaign:", campaign.Id)
		delete(Campaigns, campaign.Id)
	} else {
		log.Println("Untracking campaign clients:", campaign.Id)
		delete(campaign.Listeners, client)
	}
	CampaignsMutex.Unlock()
}

func GetCampaign(id string) *Campaign {
	var camp *Campaign
	CampaignsMutex.Lock()
	camp = Campaigns[id]
	CampaignsMutex.Unlock()
	return camp
}
func (campaign *Campaign) Broadcast(msg Message) {
	ClientsMutex.Lock()
	for client := range campaign.Listeners {
		client.Channel <- msg
	}
	ClientsMutex.Unlock()
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
