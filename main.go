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

var SoundEventQueue = make(chan SoundEvent, 1000)

type SoundEvent struct {
	// mp3/tts
	T    string
	Data string
}

var BotEventQueue = make(chan BotEvent, 100)

type BotEvent struct {
	// eve/?
	T        string
	Data     string
	Original string
}

func ProcessBotEvents() {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r)
		}
		monitor <- 12
	}()

	for s := range BotEventQueue {
		fmt.Printf("BOT QUEUE: len(%d), max(%d)", len(BotEventQueue), cap(BotEventQueue))
		switch s.T {
		case "eve":
			askTheBot(s.Data, s.Original)
		default:
			fmt.Println("UKNOWN BOT EVENT", s)
		}
	}
}

func ProcessSoundEvents() {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r)
		}
		monitor <- 11
	}()

	for s := range SoundEventQueue {
		fmt.Printf("SOUND QUEUE: len(%d), max(%d)", len(SoundEventQueue), cap(SoundEventQueue))
		switch s.T {
		case "mp3":
			PlayRewardMP3(s.Data)
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

	// go captureKeys()
	mongowrapper.InitCollections()

	InitTwitchClient()
	BaseMSG = INIT_MSG()

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
	go ProcessBotEvents()

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
			} else if ID == 12 {
				go ProcessBotEvents()
			} else if ID == 13 {
				// go captureKeys()
			}

		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}
