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

// channel / list of users
var Usermap = make(map[string][]string)
var UsermapTotal = make(map[string]int)
var UserChannelCount = make(map[string]int)

func InitUserMap() {
	Usermap["zhuffles"] = []string{}
	// Usermap["hofkari"] = []string{}
	// Usermap["kevinramm"] = []string{}
	// Usermap["justruss"] = []string{}
	// Usermap["mmarkers"] = []string{}
}
func main() {
	InitUserMap()
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

	// ScrapeTwitch()
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

		for {
			time.Sleep(time.Minute * 1)
		}
	}()

	go func() {
		// TWITCHclient = twitch.NewClient("ZENDROIDlive", os.Getenv("TWITCH_KEY"))
		TWITCHclient = twitch.NewClient("zkynettest", os.Getenv("TWITCH_KEY"))
		TWITCHclient.OnPrivateMessage(twitchMessageHandler)
		TWITCHclient.OnUserJoinMessage(twitchJoinHandler)

		err = TWITCHclient.Connect()
		if err != nil {
			panic(err)
		}

	}()

	time.Sleep(1 * time.Second)
	for i := range Usermap {
		time.Sleep(100 * time.Millisecond)

		log.Println("joining .. ", i)

		TWITCHclient.Join(i)
	}

	log.Println("sleeping for 30")
	// time.Sleep(30 * time.Second)

	// for {
	// 	log.Println("starting user map parsing..")
	// 	for i := range Usermap {
	// 		time.Sleep(300 * time.Millisecond)
	// 		users, err := TWITCHclient.Userlist(i)
	// 		if err != nil {
	// 			log.Println(err, string(debug.Stack()))
	// 		}

	// 		Usermap[i] = users
	// 		log.Println("api vs total:", i, len(users), UsermapTotal[i])
	// 	}
	// 	// log.Println("USERS:", Usermap)

	// 	for _, v := range Usermap {
	// 		for _, vv := range v {
	// 			UserChannelCount[vv]++
	// 		}
	// 	}
	// 	log.Println("--------------------------------------------------")
	// 	log.Println("total users:", len(UserChannelCount))
	// 	for i, v := range UserChannelCount {
	// 		if v > 4 {
	// 			log.Println(i, v)
	// 		}
	// 	}
	// 	reader := bufio.NewReader(os.Stdin)
	// 	_, _ = reader.ReadString('\n')

	// }
	// temp

	// Start discord client

	// SendHello()
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

	log.Println("COMMAND:", splitMessage[0], " //  KEY NR:", keyboardCommandsToActions[splitMessage[0]], " // CTL NR:", controlValue, " // TIMES:", totalButtonPresses)

	// if we don't find any controls, go back
	if keyboardCommandsToActions[splitMessage[0]] == 0 {
		return
	}

	if strings.Contains(message, "!volumeup") {
		totalButtonPresses = 10
	}

	if strings.Contains(message, "!volumedown") {
		totalButtonPresses = 10

	}
	// if strings.Contains(message, "!d") || strings.Contains(message, "!a") {
	// 	KEYBONDING.SetKeys(keyboardCommandsToActions[splitMessage[0]])
	// 	PressKey(controlValue * 8)
	// 	return
	// }

	for buttonPress := 0; buttonPress < totalButtonPresses; buttonPress++ {
		KEYBONDING.SetKeys(keyboardCommandsToActions[splitMessage[0]])
		PressKey(controlValue * oneMeter)
	}

}
func TriggerSocials(message,username,platform string){
	splitMessage := strings.Split(message, " ")
	text := socialCommandsToText[splitMessage[0]]
	if text == "" {
		return
	}

	log.Println("SOCIAL", username, platform, text)
	if platform == "twitch" {
	STREAMSendSystemMessage(TWITCHCHANNEL, text)
	TWITCHWhisperUser(username, text)
	TWITCHSendMsgToChannel(TWITCHCHANNEL, text)
	} else if platform == "discord" {
		DISCORDSendMsgToChannel(DISCORDCHANNEL, text)
	}
} 

func twitchMessageHandler(message twitch.PrivateMessage) {
	// fmt.Println(message)
	log.Println("MESSAGE:",message.Channel, message.Message)
	for _, v := range WSMAP {
		websocket.Message.Send(v, message.User.DisplayName+":xx:"+message.Message)
	}
	RaiderIOCheck(message.Message, "twitch")
	TriggerControls(message.Message)
	TriggerSocials(message.Message,  message.User.Name, "twitch")
}
func RaiderIOCheck(message string, platform string) {


	if !strings.Contains(message, "!player") {
      return
	}
	info := strings.Split(message, " ")
	Player := RaiderIOCharacter(info[1], info[2], info[3], []string{"gear"})
	log.Println(Player.Base.Gear.Items.Neck)
 
	x := []string{}
	x = append(x, "https://raider.io/characters/"+info[1]+"/"+info[2]+"/"+info[3]+"")
	if info[4] == "ilvl" {
		x = append(x, "ILVL: "+strconv.Itoa(Player.Base.Gear.ItemLevelEquipped))
		x = append(x, "Total ILVL: "+strconv.Itoa(Player.Base.Gear.ItemLevelTotal))
	} else if info[4] == "corruption" {
		x = append(x,"Equiped Corruption: "+strconv.Itoa(Player.Base.Gear.Corruption.Added))
		x = append(x,"Resisted Corruption: "+strconv.Itoa(Player.Base.Gear.Corruption.Resisted))
		x = append(x,"Total Corruption: "+strconv.Itoa(Player.Base.Gear.Corruption.Total))
	} else if info[4] == "essence" {
		for _,v := range Player.Base.Gear.Items.Neck.HeartOfAzeroth.Essences {
			log.Println("adding one essence ...")
			x = append(x,"Essence: "+v.Power.Essence.Name +" Rank("+strconv.Itoa(v.Rank)+")")
		}
	} else if info[4] == "items" {
		staticlink := ""
		if platform == "twitch" {
			staticlink = "wowhead.com/item="
		} else {
			staticlink = "https://wowhead.com/item="
		}
		x = append(x, "Head "+staticlink+ strconv.Itoa(Player.Base.Gear.Items.Head.ItemID))
		x = append(x, "Neck "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Neck.ItemID))
		x = append(x, "Shoulder "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Shoulder.ItemID))
		x = append(x, "Chest "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Chest.ItemID))
		x = append(x, "Wrist "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Wrist.ItemID))
		x = append(x, "Hands "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Hands.ItemID))
		x = append(x, "Waist "+strconv.Itoa(Player.Base.Gear.Items.Waist.ItemID))
		x = append(x, "Legs "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Legs.ItemID))
		x = append(x, "Feet "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Feet.ItemID))
		x = append(x, "Trinket "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Trinket1.ItemID))
		x = append(x, "Trinket "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Trinket2.ItemID))
		x = append(x, "Main Hand "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Mainhand.ItemID))
		x = append(x, "Off Hand "+staticlink+strconv.Itoa(Player.Base.Gear.Items.Offhand.ItemID))
	} else if info[4] == "raiderio" {
		for _,v := range Player.Base.Score {
		x = append(x, "RaiderIO dps: "+fmt.Sprintf("%.1f",v.Scores.Dps ) )	
		x = append(x, "RaiderIO healer: "+fmt.Sprintf("%.1f",v.Scores.Healer ) )	
		}
	}
 
	 time.Sleep(200*time.Millisecond)
	if platform == "twitch" {
		log.Println("Sending to twitch:",strings.Join(x, " / "))
		TWITCHSendMsgToChannel(TWITCHCHANNEL, strings.Join(x, " / "))
	} else if platform == "discord" {
		DISCORDSendMsgToChannel(DISCORDCHANNEL, strings.Join(x, " / "))
	}


}
func twitchJoinHandler(message twitch.UserJoinMessage) {
	// fmt.Println(message)
	log.Println("JOIN:", message)
	// for _, v := range WSMAP {
	// 	websocket.Message.Send(v, message.User.DisplayName+":xx:"+message.Message)
	// }
}

func discordMessageHandler(session disgord.Session, evt *disgord.MessageCreate) {
	if disgord.Snowflake(DISCORDCHANNEL) != evt.Message.ChannelID {
		log.Println("NMO MATCH:", DISCORDCHANNEL, evt.Message.ChannelID)
		return	
	}
	for _, v := range WSMAP {
		websocket.Message.Send(v, evt.Message.Author.Username+":xx:"+evt.Message.Content)
	}
	log.Println("DISCORD MESSAGE", evt.Message.ChannelID, evt.Message.Content)
		RaiderIOCheck(evt.Message.Content, "discord")
	TriggerControls(evt.Message.Content)
	TriggerSocials( evt.Message.Content,  evt.Message.Author.Username, "discord")
}


func STREAMSendSystemMessage(channel string, msg string) {
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



func DISCORDSendMsgToChannel(channel uint64, msg string) {
	_, err := DISCORDclient.CreateMessage(context.Background(), disgord.Snowflake(channel), &disgord.CreateMessageParams{
		Content: msg,
	}, disgord.Flag(1<<4))
	if err != nil {
		log.Println(err, string(debug.Stack()))
	}
}
func TWITCHSendMsgToChannel(channel, msg string) {
	TWITCHclient.Say(channel, msg)
}
func TWITCHWhisperUser(user, msg string) {
	TWITCHclient.Whisper(user, msg)
}

func SendHello(){
	_, err := DISCORDclient.CreateMessage(context.Background(), disgord.Snowflake(DISCORDCHANNEL), &disgord.CreateMessageParams{
			Content: "ZENDROID is now LIVE on twitch: https://www.twitch.tv/zendroidlive",
		}, disgord.Flag(1<<4))
		if err != nil {
			log.Println(err, string(debug.Stack()))
		}

		
		TWITCHclient.Say("zhuffles", "ZENDROID is now LIVE, follow me on twitter: https://www.twitter.com/zkynetio")
}