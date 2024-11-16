package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	tirc "github.com/gempir/go-twitch-irc/v4"
	"github.com/google/uuid"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"github.com/nicklaw5/helix"
)

type IRC_CLIENT struct {
	Name   string
	Client *tirc.Client
	// 1 == main channel
	// 2 == sub channel
	ChannelMap map[string]*IRC_CHANNEL
}

type IRC_CHANNEL struct {
	Name          string
	BroadCasterID string
	Type          int
}

func (C *IRC_CLIENT) ReplyToUser(userName string, msg string, channel string) {
	if channel == "" {
		C.Client.Say(C.Name, userName+" >> "+msg)
	} else {
		C.Client.Say(channel, userName+" >> "+msg)
	}
}

func (C *IRC_CLIENT) Reply(msg string, channel string) {
	if channel == "" {
		C.Client.Say(C.Name, msg)
	} else {
		C.Client.Say(channel, msg)
	}
}

func (C *IRC_CLIENT) GetAllChannelEmotes() {
	for _, v := range C.ChannelMap {
		log.Println("GETTING EMOTED FOR CHANNEL: ", v.Name)
		GetChannelEmotes(v.BroadCasterID)
	}
}

func NewUserSubReSubRaidMessage(user tirc.UserNoticeMessage) {
	TWITCH_CLIENT.ReplyToUser(user.User.DisplayName, "thank you!", "")
}

func USER_TEST(user tirc.UserStateMessage) {
	log.Println("State:", user)
}

func (C *IRC_CLIENT) POST_INFO() {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r, string(debug.Stack()))
		}
		monitor <- 1337
	}()
	time.Sleep(1 * time.Hour)

	if C.Client != nil {

		returnText, ok := TextCommands["!x"]
		if ok {
			TWITCH_CLIENT.Reply(returnText, "")
		}
		time.Sleep(30 * time.Second)
		returnText, ok = TextCommands["!discord"]
		if ok {
			TWITCH_CLIENT.Reply(returnText, "")
		}
		time.Sleep(30 * time.Second)
		returnText, ok = TextCommands["!vpn"]
		if ok {
			TWITCH_CLIENT.Reply(returnText, "")
		}
	}
}

func (C *IRC_CLIENT) Connect() {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r, string(debug.Stack()))
		}
		monitor <- 10
	}()

	log.Println("KEY LENGTH: ", len(os.Getenv("TWITCH_KEY")))
	C.Client = tirc.NewClient(C.Name, os.Getenv("TWITCH_KEY"))
	C.Client.SendPings = true
	C.Client.IdlePingInterval = time.Duration(time.Second * 10)
	C.Client.PongTimeout = time.Duration(time.Second * 60)
	C.Client.OnPrivateMessage(NewMessage)
	C.Client.OnUserNoticeMessage(NewUserSubReSubRaidMessage)
	// C.Client.OnUserStateMessage(USER_TEST)

	go func() {
		time.Sleep(3 * time.Second)
		C.JoinChannels()
		time.Sleep(3 * time.Second)
		go TWITCH_CLIENT.POST_INFO()
		TWITCH_CLIENT.Reply("BOT ONLINE!", "")
	}()

	err := C.Client.Connect()
	if err != nil {
		log.Println(err)
	}
}

func PlaceEventOnSoundQueue(t string, data string) {
	select {
	case SoundQueue <- SoundEvent{
		T:    t,
		Data: data,
	}:
	default:
		fmt.Println("SOUND QUEUE FULL")
		fmt.Println("SOUND QUEUE FULL")
		fmt.Println("SOUND QUEUE FULL")
		fmt.Println("SOUND QUEUE FULL")
	}
	return
}

func (C *IRC_CLIENT) JoinChannels() {
	for _, v := range C.ChannelMap {
		log.Println("JOINING CHANNEL: ", v)
		C.Client.Join(v.Name)
		GetChannelEmotes(v.BroadCasterID)
	}
}

func NewMessage(msg tirc.PrivateMessage) {
	ProcessMessage(msg)
}

func ProcessMessage(msg tirc.PrivateMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println(string(debug.Stack()))
		}
	}()

	// lowerUser := strings.ToLower(msg.User.DisplayName)
	// fmt.Println("---------------------")
	// fmt.Println(msg.FirstMessage)
	// fmt.Println(msg.Action)
	// fmt.Println(msg.Bits)
	// fmt.Println(msg.Tags)
	// fmt.Println(msg.Type)
	// for _, v := range msg.Emotes {
	// 	fmt.Println(v.ID, v.Name, v.Count)
	// 	for _, vv := range v.Positions {
	// 		fmt.Println(vv.Start, vv.End)
	// 	}
	// }
	// fmt.Println(msg.CustomRewardID)
	// fmt.Println(msg.User.ID)
	// fmt.Println(msg.User.Name)
	// fmt.Println(lowerUser)
	// fmt.Println(msg.Channel, msg.Message)
	// fmt.Println("---------------------")

	U, err := FindOrUpsertUser(&msg.User)
	if err != nil {
		U.DisplayName = msg.User.DisplayName
		U.Name = msg.User.Name
		U.ID = msg.User.ID
		U.Color = msg.User.Color
		U.Badges = msg.User.Badges
		err = IncrementUserPoints(U, 500)
		if err == nil {
			U.Points = 100
		}
	} else {
		_ = IncrementUserPoints(U, 5)
	}

	log.Println("CUSTOM REWARD ID:", msg.CustomRewardID)
	returnText, ok := TextCommands[msg.Message]
	if ok {
		TWITCH_CLIENT.Reply(returnText, "")
		return
	}

	if CheckCustomReward(U, msg) {
		return
	}

	mp3, ok := MP3Map[msg.CustomRewardID]
	if ok {
		go PlaceEventOnSoundQueue("mp3", mp3)
		// go PlayMP3(mp3)
		return
	}

	if strings.Contains(msg.Message, "!time") {
		TWITCH_CLIENT.Reply(time.Now().Format(time.RFC3339), "")
		return
	}

	if strings.Contains(msg.Message, "!bot") {
		if msg.User.DisplayName == "KEYB1ND_" {
			go askTheBot(msg.Message)
		}
		return
	}

	// banword := ""
	// isBanned := false
	// if strings.Contains(msg.Message, " tailwind ") || strings.HasPrefix(msg.Message, "tailwind") {
	// 	banword = "tailwind"
	// 	isBanned = true
	// }

	// if strings.Contains(msg.Message, " rust ") || strings.HasPrefix(msg.Message, "rust") {
	// 	banword = "rust"
	// 	isBanned = true
	// }

	// if isBanned {
	// 	TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "You said a naughty word: "+banword+" -1000 points for you.", "")
	// 	_ = IncrementUserPoints(U, -1000)
	// 	return
	// }

	// if strings.Contains(msg.Message, "!tts") {
	// 	go CustomTTS(*U, msg)
	// 	return
	// }

	if strings.Contains(msg.Message, "!quote help") || msg.Message == "!quote" {
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "To make a quote use 'Quote:' before your sentence...... example: 'Quote: This is a quote!' ", "")
		return
	}

	if strings.Contains(msg.Message, "!points") {
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "you have >> "+strconv.Itoa(U.Points)+" Points", "")
		return
	}

	if strings.Contains(msg.Message, "!quote") {
		RandQuote(&msg)
		return
	}

	if strings.Contains(msg.Message, "!top10") {
		Top10Command()
		return
	}

	if strings.Contains(msg.Message, "!roll") {
		err := UserRollCommand(U, &msg)
		if err != nil {
			log.Println("Error While Rolling: ", err)
			return
		}
		return
	}

	_ = SaveMessage(&msg)
	// log.Println("USER FROM DB: ", U)
}

var (
	RollTimeout = make(map[string]time.Time)
	RollLock    sync.Mutex
)

func RandQuote(msg *tirc.PrivateMessage) {
	splitMatch := strings.Split(msg.Message, " ")
	if len(splitMatch) < 2 || len(splitMatch) > 2 {
		return
	}

	userToLower := strings.ToLower(splitMatch[1])
	userToLower = strings.Replace(userToLower, "@", "", -1)
	msgs, err := FindUserMessagesFromMatch(userToLower, "Quote:")
	if err != nil {
		return
	}
	if len(msgs) == 0 {
		TWITCH_CLIENT.Reply("No qoutes found for "+userToLower, "")
		return

	}

	random := rand.Intn(len(msgs))
	selectedMSG := msgs[random]

	outMSG := strings.Replace(selectedMSG.Message, "Quote:", "", -1)
	go PlayTTS(outMSG)
	TWITCH_CLIENT.Reply(selectedMSG.User.DisplayName+" '' "+outMSG+" '' - "+selectedMSG.Time.Format("Mon 02 Jan 15:04:05 MST 2006"), "")
}

func PlayMP3(tag string) {
	// cmd := exec.Command("ffplay", "-v", "0", "-nodisp", "-autoexit", "./mp3/"+tag+".mp3")
	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Println(err, string(out))
	// }
	f, err := os.Open("./mp3/" + tag + ".mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func Top10Command() {
	userList := GetTop10()
	outMsg := ""
	rank := 1
	for _, v := range userList {
		if v.Name == USERNAME {
			continue
		}
		outMsg += strconv.Itoa(rank) + "#" + v.DisplayName + "(" + strconv.Itoa(v.Points) + ") ......."
		rank++
	}
	TWITCH_CLIENT.Reply(outMsg, "")
}

func UserRollCommand(user *User, msg *tirc.PrivateMessage) (err error) {
	rollSplit := strings.Split(msg.Message, " ")
	if len(rollSplit) < 2 {
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "Invalid roll format", "")
		return
	}

	rollAmount, err := strconv.Atoi(rollSplit[1])
	if err != nil {
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "Invalid roll format", "")
		return
	}

	if user.Points < rollAmount || rollAmount < 0 {
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "you do not have enough points to gamba", "")
		return
	}

	lastRoll, ok := RollTimeout[msg.User.ID]
	if ok {
		seconds := time.Since(lastRoll).Seconds()
		if seconds < 20 {
			TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, strconv.Itoa(int(20-seconds))+" seconds until you can roll again", "")
			return
		}
	}
	RollTimeout[msg.User.ID] = time.Now()

	random := rand.Intn(100) + 1
	if random < 30 {
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, " rolls "+strconv.Itoa(random)+" and wins nothing", "")
		_ = IncrementUserPoints(user, -rollAmount)

	} else if random%2 == 0 {

		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, " rolls "+strconv.Itoa(random)+" and wins "+strconv.Itoa(rollAmount), "")
		_ = IncrementUserPoints(user, rollAmount)

	} else {

		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, " rolls "+strconv.Itoa(random)+" and wins nothing", "")
		_ = IncrementUserPoints(user, -rollAmount)
	}

	return
}

func CreateAPIClient() {
	twitchkey := os.Getenv("TWITCH_KEY")
	twitchkey = strings.Split(twitchkey, ":")[1]

	var err error
	HELIX_CLIENT, err = helix.NewClient(&helix.Options{
		AppAccessToken: twitchkey,
		ClientSecret:   os.Getenv("CLIENT_SECRET"),
		ClientID:       os.Getenv("CLIENT_ID"),
	})
	if err != nil {
		log.Println("UNABLE TO CREATE API CLIENT: ", err)
	}
}

func GetGlobalEmotes() {
	resp, err := HELIX_CLIENT.GetGlobalEmotes()
	if err != nil {
		log.Println("ERROR GETTING GLOBAL EMOTES:", err)
		return
	}

	for _, v := range resp.Data.Emotes {
		EmoteMap[v.ID] = v
	}
}

func GetChannelEmotes(BroadCasterID string) {
	resp, err := HELIX_CLIENT.GetChannelEmotes(&helix.GetChannelEmotesParams{
		BroadcasterID: BroadCasterID,
	})
	if err != nil {
		log.Println("ERROR GETTING CHANNEL EMOTES:", err)
		return
	}

	for _, v := range resp.Data.Emotes {
		EmoteMap[v.ID] = v
	}
}

// https://github.com/coqui-ai/TTS
// https://github.com/coqui-ai/TTS
// python3 TTS/server/server.py --vocoder_name vocoder_models/en/ljspeech/hifigan_v2
// jenny is good too
func CustomTTS(user User, msg tirc.PrivateMessage) {
	if user.Points < 50 {
		TWITCH_CLIENT.ReplyToUser(user.DisplayName, "You need 200 points for TTS", "")
		return
	}

	splitTTS := strings.Split(msg.Message, "tts")
	if len(splitTTS) < 2 {
		TWITCH_CLIENT.ReplyToUser(user.DisplayName, "TTS was badly formatted.. try: !tts [MSG]", "")
		return
	}

	err := IncrementUserPoints(&user, -50)
	if err != nil {
		TWITCH_CLIENT.ReplyToUser(user.DisplayName, "We could not play your sound clip!", "")
		return
	}

	PlayTTS(splitTTS[1])
}

func PlayTTS(msg string) {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r, string(debug.Stack()))
		}
	}()

	// http://127.0.0.1:5002/api/tts?text=whats%20up&speaker_id=&style_wav=&language_id=

	params := url.Values{}
	params.Add("text", msg)

	httpClient := new(http.Client)
	// urlmsg := url.QueryEscape(msg)
	req, err := http.NewRequest("GET", "http://127.0.0.1:5002/api/tts?"+params.Encode(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	fileName := uuid.NewString() + "-" + strconv.Itoa(int(time.Now().UnixNano()))

	f, err := os.Create("./tts/" + fileName + ".wav")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	f2, err := os.Create("./tts/" + fileName + ".txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer f2.Close()
	_, err = f.Write(bytes)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = f2.Write([]byte(msg))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("BYTES:", len(bytes))

	PlayTTSFile(f.Name())
}

func PlayTTSFile(tag string) {
	f, err := os.Open(tag)
	if err != nil {
		log.Println("error opening tts file:", err)
		return
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

type BotModel struct {
	Model     string    `json:"model"`
	Prompt    string    `json:"prompt"`
	Stream    bool      `json:"stream"`
	Options   *BotOpts  `json:"options"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []*BotMSG `json:"messages"`
}

type BotMSG struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BotOpts struct {
	Temperature float64 `json:"temperature"`
	// TopP        float64 `json:"top_p"`
	// TopK        int     `json:"top_k"`
	// MinP        float64 `json:"min_p"`
	// NumPredict  int     `json:"num_predict"`
}

type BotResp struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

type BotRespChat struct {
	Message *BotMSG `json:"message"`
}

var BaseMSG = make([]*BotMSG, 0)

func INIT_MSG() {
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "user",
		Content: "I want you to refer to me as Key from now on.",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "assistant",
		Content: "I will refer to you as key from now on.",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "user",
		Content: "I want you to pretend to be evil, be subtle about it but try to be a bit snarky with your replies. Also, try to make fun of me in subtle ways when possible, don't over do it though.",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "assistant",
		Content: "I will do that",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "user",
		Content: "Please try to keep your answers short.",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "assistant",
		Content: "I will do that",
	})
}

func askTheBot(msg string) {
	ms := strings.Split(msg, " ")
	m := strings.Join(ms[1:], " ")

	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "user",
		Content: m,
	})

	BM := new(BotModel)
	BM.Model = "mannix/llama3.1-8b-abliterated"
	// BM.Prompt = m
	BM.Messages = BaseMSG
	BM.Stream = false
	BM.Options = new(BotOpts)
	BM.Options.Temperature = 0.8
	BM.MaxTokens = 50

	ob, err := json.Marshal(BM)

	buff := bytes.NewBuffer(ob)

	// req, err := http.NewRequest("POST", "http://localhost:11434/api/generate", buff)
	httpClient := new(http.Client)
	req, err := http.NewRequest("POST", "http://localhost:11434/api/chat", buff)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	BR := new(BotRespChat)
	err = json.Unmarshal(bytes, BR)
	if err != nil {
		fmt.Println(string(bytes))
		fmt.Println("lama resp err:", err)
		return
	}

	BaseMSG = append(BaseMSG, BR.Message)

	fmt.Println(BR.Message.Content)
	PlayTTS(BR.Message.Content)
}
