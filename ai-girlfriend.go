package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

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
	// BaseMSG = append(BaseMSG, &BotMSG{
	// 	Role:    "user",
	// 	Content: "I want you to pretend to be evil, be subtle about it but try to be a bit snarky with your replies. Also, try to make fun of me in subtle ways when possible, don't over do it though.",
	// })
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "user",
		Content: "I want you to be a bit evil and snarky. But in a funny way. It's ok to be a bit mean.",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "assistant",
		Content: "I will do that",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "user",
		Content: "I want you to sometimes answer my questions with your own questions",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "assistant",
		Content: "I will do that",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "user",
		Content: "Please try to keep your answers below or around 100 words.",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "assistant",
		Content: "I will do that",
	})
	BaseMSG = append(BaseMSG, &BotMSG{
		Role:    "user",
		Content: "I want you to refer to yourself as Vespera and your visual avatar is a gothic chick.",
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
