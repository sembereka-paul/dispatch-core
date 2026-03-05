package main

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"
)

func TestSSEConnectionOk(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/notifications/new", nil)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "text/event-stream")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	defer resp.Body.Close()
}
