package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

type OAUTHS struct {
	X     string `json:"type"`
	Nonce string
	Data  Data
}
type Data struct {
	Topics     []string
	Auth_token string
}

var (
	clientID = "jdjram29pjgewueoekjq0pmwax8ace"
	// Consider storing the secret in an environment variable or a dedicated storage system.
	clientSecret = "o1665kvpw6eom5z1xqp9tx6an6jm8n"
	channelID    = "529902238"
	oauth2Config *clientcredentials.Config
	T            = "x4vrec5pyspt4peawmwkmf8ef3moer"
	UT           = "3zn08h32eahob5jy8qyekrhjri5iy6"
)

func serverFS() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	log.Println("Listening on :80...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {

	if T == "" {

		oauth2Config = &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     twitch.Endpoint.TokenURL,
			Scopes:       []string{"bits:read,whispers:edit"},
		}

		token, err := oauth2Config.Token(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Access token: %s %s\n", token.AccessToken, token.TokenType)
		T = token.AccessToken
	}

	log.Println("Using token", T)
	// getUserInfo(T)
	// getFollows(T)
	// getSubs(T)
	go serverFS()
	// go launchWS()

	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for {
		s := <-signal_chan
		switch s {
		// kill -SIGHUP XXXX
		case syscall.SIGHUP:
			fmt.Println("hungup")
			os.Exit(1)
		// kill -SIGINT XXXX or Ctrl+c
		case syscall.SIGINT:
			fmt.Println("Warikomi")

			os.Exit(1)
		// kill -SIGTERM XXXX
		case syscall.SIGTERM:
			fmt.Println("force stop")
			os.Exit(1)

		// kill -SIGQUIT XXXX
		case syscall.SIGQUIT:
			fmt.Println("stop and core dump")
			os.Exit(1)

		default:
			fmt.Println("Unknown signal.")
			os.Exit(1)
		}
	}

}

func getUserInfo(token string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?login=zendroidlive", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Client-ID", clientID)
	resp, err := client.Do(req)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Println(bodyString)

	defer resp.Body.Close()
}
func getSubs(token string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/webhooks/subscriptions", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Client-ID", clientID)
	resp, err := client.Do(req)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Println(bodyString)

	defer resp.Body.Close()
}
func getFollows(token string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users/follows?to_id="+channelID, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Client-ID", clientID)
	resp, err := client.Do(req)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	bodyString := string(bodyBytes)
	log.Println(bodyString)

	defer resp.Body.Close()
}

func getBits(token string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/bits/leaderboard?count=10", nil)
	if err != nil {
		panic(err)
	}
	log.Println("inner token", token)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Client-ID", clientID)
	resp, err := client.Do(req)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Println(bodyString)

	defer resp.Body.Close()
}

func launchWS() {

	userOauth := OAUTHS{
		X:     "LISTEN",
		Nonce: "q1NhD7JfiJ0G6Zx",
		Data: Data{
			Topics:     []string{"channel-subscribe-events-v1." + channelID, "whispers." + channelID, "chat_moderator_actions." + channelID, "channel-bits-events-v2." + channelID},
			Auth_token: UT,
		},
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial("wss://pubsub-edge.twitch.tv", nil)
	if err != nil {
		log.Println("dial:", err)
		return
	}
	defer c.Close()

	mwg, err := json.Marshal(userOauth)
	if err != nil {
		log.Println(err, string(debug.Stack()))
	}
	log.Println(string(mwg))
	err = c.WriteMessage(websocket.TextMessage, mwg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	log.Println("reading ...")
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", string(message))
	}

}
