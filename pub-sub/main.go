package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	pb "coop/proto"

	"google.golang.org/grpc"
)

func parseLine(line string) (string, string) {
	if i := strings.IndexByte(line, ':'); i >= 0 {
		key := strings.TrimSpace(line[:i])
		val := line[i+1:]
		val = strings.TrimRight(val, "\r\n")
		if key == "" {
			return "", ""
		}
		return key, val
	}
	return "", ""
}

type Message struct {
	event string `json:"event"`
	data  string `json:"data"`
	tag   string `json:"tag"`
}

type server struct {
	pb.UnimplementedEventServer
}

var tag = make(chan string, 1)
var out = make(chan pb.EventReply)

func (s *server) Sub(in *pb.SubscribeRequest, stream pb.Event_SubServer) error {
	select {
	case tag <- in.Tag:
	default:
		return errors.New("Failed to sub")
	}

	for res := range out {
		stream.Send(&res)
	}
	return nil
}

func processLine(scanner *bufio.Scanner) (string, string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", "", err
		}
	}
	first := scanner.Text()
	event, data := parseLine(first)
	return event, data, nil
}

var (
	port = flag.Int("port", 50052, "The server port")
)
var (
	MASTODON_BASE_URL     = os.Getenv("MASTODON_BASE_URL")
	MASTODON_ACCESS_TOKEN = os.Getenv("MASTODON_ACCESS_TOKEN")
)

// Managens a single subscription to ta tag
// take the subscription name as sub and an out channel to write to.
func subscriptionWorker(sub string, out chan<- pb.EventReply) {
	req, err := http.NewRequest("GET", MASTODON_BASE_URL+"/api/v1/streaming/hashtag?tag="+sub, nil)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+MASTODON_ACCESS_TOKEN)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error")
		return
	}

	scanner := bufio.NewScanner(resp.Body)

	// attempt reading pairs
	// first find key, then look for value
	next := ""
	msg := Message{}
	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Println("Error reading line", err)
				break
			}
		}
		first := scanner.Text()
		field, data := parseLine(first)
		if next == "" && field == "event" {
			fmt.Println(field, data)
			msg.event = data
			next = "data"
		} else if next == "data" && field == "data" {
			fmt.Println(field, data)
			msg.data = data
			next = ""
		}

		if msg.data != "" {
			msg.tag = sub

			out <- pb.EventReply{
				Event: msg.event,
				Data:  msg.data,
				Tag:   msg.tag,
			}
			msg = Message{}
		}
	}
}

func main() {
	defer close(tag)
	defer close(out)

	go func() {
		for sub := range tag {
			go subscriptionWorker(sub, out)
		}
	}()

	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svr := grpc.NewServer()
	pb.RegisterEventServer(svr, &server{})
	log.Printf("server listening at %v", listener.Addr())

	if err := svr.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
