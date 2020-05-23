package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/gempir/go-twitch-irc"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/micmonay/keybd_event"
	"golang.org/x/net/websocket"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	LoadMaps()

	KEYBONDING, err = keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	// Start twitch client
	go func() {
		TWITCHclient = twitch.NewClient("zendroidlive", os.Getenv("TWITCH_KEY"))
		TWITCHclient.Join("zendroidlive")
		TWITCHclient.OnPrivateMessage(twitchMessageHandler)

		err = TWITCHclient.Connect()
		if err != nil {
			panic(err)
		}
	}()

	// Start discord client
	go func() {
		DISCORDclient = disgord.New(disgord.Config{
			BotToken: os.Getenv("DISCORD_KEY"),
		})

		if err = DISCORDclient.Connect(context.Background()); err != nil {
			log.Println(err)
			return
		}

		defer DISCORDclient.Disconnect()
		log.Println("starting discord listener.")
		DISCORDclient.On(disgord.EvtMessageCreate, discordMessageHandler)
		// TODO: implement a watcher
		for {
			time.Sleep(time.Minute * 1)
		}
	}()

	e := echo.New()
	e.Static("/", "./ws.html")
	e.GET("/ws", hello)
	e.Start(":1234")
}

func PressKey(controlValue int) {
	err := KEYBONDING.Press()
	if err != nil {
		log.Println("KEYPRESSERR:", err)
		return
	}
	time.Sleep(time.Duration(controlValue) * time.Millisecond)
	err = KEYBONDING.Release()
	if err != nil {
		log.Println("KEYPRESSERR:", err)
		return
	}
}

// 500 ...
// 45 daagrees in the turn
// 2 meters foward
// 1 meter backwards
func TriggerControls(message string) {
	var controlValue = 0
	var err error
	var totalButtonPresses = 1
	splitMessage := strings.Split(message, " ")
	if len(splitMessage) > 1 {
		controlValue, err = strconv.Atoi(splitMessage[1])
		if err != nil {
			log.Println(err)
			return
		}
	}

	// if we don't find any controls, go back
	if keyboardCommandsToActions[splitMessage[0]] == 0 {
		return
	}

	if strings.Contains("!volumeup", message) {
		totalButtonPresses = 10
	}

	if strings.Contains("!volumedown", message) {
		totalButtonPresses = 10
	}

	log.Println("COMMAND:", splitMessage[0], " //  KEY NR:", keyboardCommandsToActions[splitMessage[0]], " // CTL NR:", controlValue, " // TIMES:", totalButtonPresses)
	for buttonPress := 0; buttonPress < totalButtonPresses; buttonPress++ {
		KEYBONDING.SetKeys(keyboardCommandsToActions[splitMessage[0]])
		PressKey(controlValue * oneMeter)
	}

}

func TriggerDiscordSocials(kb keybd_event.KeyBonding, message string) {
	splitMessage := strings.Split(message, " ")
	text := socialCommandsToText[splitMessage[0]]
	if text == "" {
		return
	}
	_, err := DISCORDclient.CreateMessage(context.Background(), disgord.Snowflake(710530877142335651), &disgord.CreateMessageParams{
		Content: text,
	}, disgord.Flag(1<<4))
	if err != nil {
		log.Println(err, string(debug.Stack()))
	}

}
func TriggerTwitchSocials(message twitch.PrivateMessage) {
	splitMessage := strings.Split(message.Message, " ")
	text := socialCommandsToText[splitMessage[0]]
	if text == "" {
		return
	}
	SendCustomSystemMessage("ZENDROIDlive", text)
	SendCustomWhisper(message.User.Name, text)
}
func twitchMessageHandler(message twitch.PrivateMessage) {
	// fmt.Println(message)
	log.Println("MESSAGE:", message.Message)
	for _, v := range WSMAP {
		websocket.Message.Send(v, message.User.DisplayName+":xx:"+message.Message)
	}
	TriggerControls(message.Message)
	TriggerTwitchSocials(message)
}

// MSG FORMAT: zkynet#5018{670416323792207872}: test message to bot
func discordMessageHandler(session disgord.Session, evt *disgord.MessageCreate) {
	for _, v := range WSMAP {
		websocket.Message.Send(v, evt.Message.Author.Username+":xx:"+evt.Message.Content)
	}
	TriggerControls(evt.Message.Content)
	TriggerDiscordSocials(KEYBONDING, evt.Message.Content)
}

func SendCustomWhisper(user string, msg string) {
	TWITCHclient.Whisper(user, msg)
}
func SendCustomSystemMessage(channel string, msg string) {
	TWITCHclient.Say(channel, msg)
	for _, v := range WSMAP {
		websocket.Message.Send(v, "system::"+msg)
	}
}

func hello(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		log.Println("new client", ws.Request().RemoteAddr)
		WSMAP[ws.Request().RemoteAddr] = ws
		for {

			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				return
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
