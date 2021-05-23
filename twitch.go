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
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

type WOWStream struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
	TagIds       []string  `json:"tag_ids"`
}
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
	clientSecret = "3nkxrurunb62tmgu8kg7bhpjh4q2i4"
	channelID    = "529902238"
	oauth2Config *clientcredentials.Config
	T            = "n91082svn543nykgd9gi9me2rc9cgg"
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
func ScrapeTwitch() {

	if T == "" {

		oauth2Config = &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     twitch.Endpoint.TokenURL,
			Scopes:       []string{"bits:read"},
		}

		token, err := oauth2Config.Token(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Access token: %s %s\n", token.AccessToken, token.TokenType)
		T = token.AccessToken
	}

	// getUserInfo(T)
	// os.Exit(1)
	// f, err := os.Create("users")
	// if err != nil {
	// 	log.Println(err, string(debug.Stack()))
	// 	panic(1)
	// }

	log.Println("Using token", T)
	streamLength := 100
	cursor := ""
	for {
		if streamLength > 1 {
			log.Println("getting channels ..", streamLength)
			time.Sleep(100 * time.Millisecond)
			Streams := getStreams(T, cursor)

			cursor = Streams.Pagination.Cursor
			streamLength = len(Streams.Data)
			for _, v := range Streams.Data {
				if v.ViewerCount > 50 && v.ViewerCount < 300 {
					Usermap[strings.ToLower(v.UserName)] = []string{}
					UsermapTotal[strings.ToLower(v.UserName)] = v.ViewerCount
				}
				// _, _ = f.WriteString(strconv.Itoa(v.ViewerCount) + " " + v.UserName + "\n")
			}
			// log.Println(streamLength, Streams.Pagination.Cursor)
		} else {
			log.Println("done getting channels..")
			return
		}
	}

}

type STREAMResponse struct {
	Data       []WOWStream
	Pagination struct {
		Cursor string
	}
}

func getStreams(token string, cursor string) *STREAMResponse {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?game_id=18122&first=100&after="+cursor, nil)
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

	var SR STREAMResponse
	err = json.Unmarshal(bodyBytes, &SR)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		return nil
	}

	defer resp.Body.Close()
	return &SR
}
func getUserInfo(token string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?login=oyd123", nil)
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
