package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/micmonay/keybd_event"
	"golang.org/x/net/websocket"
)

var WSMAP = make(map[string]*websocket.Conn)
var oneMeter = 200

// 500 ...
// 45 daagrees in the turn
// 2 meters foward
// 1 meter backwards

func PressKey(kb keybd_event.KeyBonding, controlValue int) {
	kb.Press()
	time.Sleep(time.Duration(controlValue) * time.Millisecond)
	kb.Release()
}

func CheckVolumeControls(kb keybd_event.KeyBonding, message twitch.PrivateMessage) {

	if strings.Contains("!next", message.Message) {
		kb.SetKeys(keybd_event.VK_MEDIA_NEXT_TRACK)
		PressKey(kb, 0)
		return
	}

	if strings.Contains("!prev", message.Message) {
		kb.SetKeys(keybd_event.VK_MEDIA_NEXT_TRACK)
		PressKey(kb, 0)
		return
	}

	if strings.Contains("!play", message.Message) {
		kb.SetKeys(keybd_event.VK_MEDIA_PLAY_PAUSE)
		PressKey(kb, 0)
		return
	}

	// splitMessage := strings.Split(message.Message, " ")
	// if len(splitMessage) < 2 {
	// 	return
	// }
	// controlValue, err := strconv.Atoi(splitMessage[1])
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	if strings.Contains("!volumeup", message.Message) {
		for i := 0; i <= 10; i++ {
			kb.SetKeys(keybd_event.VK_VOLUME_UP)
			PressKey(kb, 0)
		}
		return
	}

	if strings.Contains("!volumedown", message.Message) {
		for i := 0; i <= 10; i++ {
			kb.SetKeys(keybd_event.VK_VOLUME_DOWN)
			PressKey(kb, 0)
		}
		return
	}
}

func CheckingForControls(kb keybd_event.KeyBonding, message twitch.PrivateMessage) {

	if strings.Contains("!mount", message.Message) {
		kb.SetKeys(keybd_event.VK_V)
		PressKey(kb, 0)
		return
	} else if strings.Contains(message.Message, "!jump") {
		kb.SetKeys(keybd_event.VK_SPACE)
		PressKey(kb, 0)
		return
	} else if strings.Contains(message.Message, "!inventory") {
		kb.SetKeys(keybd_event.VK_B)
		PressKey(kb, 0)
		return
	} else if strings.Contains(message.Message, "!dance") {
		kb.SetKeys(keybd_event.VK_K)
		PressKey(kb, 0)
		return
	}

	splitMessage := strings.Split(message.Message, " ")
	if len(splitMessage) < 2 {
		return
	}
	controlValue, err := strconv.Atoi(splitMessage[1])
	if err != nil {
		log.Println(err)
		return
	}
	if strings.Contains(message.Message, "!w") {
		kb.SetKeys(keybd_event.VK_W)
		PressKey(kb, controlValue*oneMeter)
	} else if strings.Contains(message.Message, "!s") {
		kb.SetKeys(keybd_event.VK_S)
		PressKey(kb, controlValue*(oneMeter*2))
	} else if strings.Contains(message.Message, "!a") {
		kb.SetKeys(keybd_event.VK_A)
		PressKey(kb, controlValue*10)
	} else if strings.Contains(message.Message, "!d") {
		kb.SetKeys(keybd_event.VK_D)
		PressKey(kb, controlValue*10)
	}
}

func CheckForSocials(kb keybd_event.KeyBonding, message twitch.PrivateMessage, client *twitch.Client) {
	if strings.Contains(message.Message, "!twitter") {
		SendCustomSystemMessage("ZENDROIDlive", "https://www.twitter.com/zkynetio", client)
		SendCustomWhisper(message.User.Name, "https://www.twitter.com/zkynetio", client)
	} else if strings.Contains(message.Message, "!youtube") {
		SendCustomSystemMessage("ZENDROIDlive", "https://www.youtube.com/channel/UCW6eiMiVqYroPX1qiosAbnQ?view_as=subscriber", client)
		SendCustomWhisper(message.User.Name, "https://www.youtube.com/channel/UCW6eiMiVqYroPX1qiosAbnQ?view_as=subscriber", client)
	}
}
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Static("/", "./ws.html")
	e.GET("/ws", hello)
	go e.Start(":1234")

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	client := twitch.NewClient("zendroidlive", os.Getenv("TWITCH_KEY"))
	client.Join("zendroidlive")
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message)
		log.Println("CONTENT:", message.Message)
		for _, v := range WSMAP {
			websocket.Message.Send(v, message.User.DisplayName+":xx:"+message.Message)
		}
		CheckVolumeControls(kb, message)
		CheckingForControls(kb, message)
		CheckForSocials(kb, message, client)
	})

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}
func SendCustomWhisper(user string, msg string, client *twitch.Client) {
	client.Whisper(user, msg)
}
func SendCustomSystemMessage(channel string, msg string, client *twitch.Client) {
	client.Say(channel, msg)
	for _, v := range WSMAP {
		websocket.Message.Send(v, "system::"+msg)
	}
}

func hello(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		log.Println("new client", ws.Request().RemoteAddr)
		WSMAP[ws.Request().RemoteAddr] = ws
		// Write
		// err := websocket.Message.Send(ws, "system::Welcome to zkynets chat bot system")
		// if err != nil {
		// 	c.Logger().Error(err)
		// }

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
