package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
	"github.com/zveinn/twitch-bot/mongowrapper"
)

var (
	monitor      = make(chan int, 100)
	MAIN_CLIENT  *IRC_CLIENT
	HELIX_CLIENT *helix.Client
)

var (
	TextCommands = make(map[string]string)
	MP3Map       = make(map[string]string)
	EmoteMap     = make(map[string]helix.Emote)
)

var SoundQueue = make(chan SoundEvent, 1000)

type SoundEvent struct {
	// mp3/tts
	T    string
	Data string
}

func ProcessSoundEvents() {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r)
		}
		monitor <- 7
	}()

	for s := range SoundQueue {
		fmt.Printf("SOUND QUEUE: len(%d), max(%d)", len(SoundQueue), cap(SoundQueue))
		switch s.T {
		case "mp3":
			PlayMP3(s.Data)
		case "tts":
			PlayTTS(s.Data)
		default:
			fmt.Println("UKNOWN SOUND EVENT", s)
		}
	}
}

var TWITCH_CLIENT = new(IRC_CLIENT)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// MakeNewToken()
	// os.Exit(1)

	fmt.Println("MONGO CONNECTING:", os.Getenv("DB"))
	err = mongowrapper.Connect(os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	mongowrapper.InitCollections()

	InitTwitchClient()

	InitCommands()
	InitMP3Map()

	err = RenewTokens()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	CreateAPIClient()

	GetGlobalEmotes()
	TWITCH_CLIENT.GetAllChannelEmotes()
	go ProcessSoundEvents()

	// go RenewTokensLoop()

	go TWITCH_CLIENT.Connect()

	for {
		select {

		case ID := <-monitor:
			log.Println("ID RETURNED: ", ID)

			if ID == 10 {
				TWITCH_CLIENT.Connect()
			} else if ID == 7 {
				go RenewTokensLoop()
			} else if ID == 1337 {
				go TWITCH_CLIENT.POST_INFO()
			} else if ID == 11 {
				go ProcessSoundEvents()
			}

		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}
