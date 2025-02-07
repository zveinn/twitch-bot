// Record Windows Audio project main.go
package mic

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/gorilla/websocket"
)

func TalkToEve() (msg string) {
	os.Remove("output.mp3")

	cmd := exec.Command("./ffmpeg.exe",
		"-f", "dshow",
		"-i", `audio=Microphone (Yeti Classic)`,
		"-t", "5",
		"-acodec", "libmp3lame",
		"-ac", "1",
		"-ar", "16000",
		"output.mp3")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("capture err:", err)
		return ""
	}
	// fmt.Println(string(out))
	// fmt.Println(err)

	fmt.Println("ABOUT TO TRANSCRIBE")
	return transcribeV2llama("output.mp3")
}

func check(err error) {

	if err != nil {
		log.Println(err, string(debug.Stack()))
		// panic(err)
	}
}

const Host = "localhost"
const Port = "2700"
const buffsize = 1_000_000

type Message struct {
	Result []struct {
		Conf  float64
		End   float64
		Start float64
		Word  string
	}
	Text string
}

var m Message

func transcribe(fn string) (msg string) {

	u := url.URL{Scheme: "ws", Host: Host + ":" + Port, Path: ""}
	fmt.Println("connecting to ", u.String())

	// Opening websocket connection
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	check(err)
	defer c.Close()

	f, err := os.Open(fn)
	check(err)
	if f == nil {
		fmt.Println("NO FILE, ", fn)
		return
	}

	// streamer, format, err := wav.Decode(f)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer streamer.Close()
	// Send configuration
	config := map[string]interface{}{
		"config": map[string]interface{}{
			"sample_rate": 16000, // Assuming the audio is at 16kHz
			// "sample_rate": format.SampleRate, // Assuming the audio is at 16kHz
		},
	}
	err = c.WriteJSON(config)
	if err != nil {
		log.Fatal("write json:", err)
	}

	// for {
	buff, err := io.ReadAll(f)

	if len(buff) == 0 && err == io.EOF {
		err = c.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}"))
		check(err)
		return ""
	}
	check(err)

	err = c.WriteMessage(websocket.BinaryMessage, buff)
	check(err)

	// Read message from server
	_, x2, errx := c.ReadMessage()
	check(errx)
	fmt.Println("YOU SAID: ", string(x2))

	c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	return string(x2)
	// 	break
	// }

	// Read final message from server
	// _, msg, err := c.ReadMessage()
	// fmt.Println("OUT:", string(msg))
	// check(err)

	// Closing websocket connection
	// Unmarshalling received message
	// err = json.Unmarshal(msg, &m)
	// check(err)
	// fmt.Println(m)
}

func transcribeV2llama(fn string) (msg string) {

	// Read the file content
	fileBytes, err := os.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}

	// Make an HTTP POST request
	resp, err := http.Post("http://localhost:8080/upload", "application/octet-stream", bytes.NewReader(fileBytes))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
	return string(body)
}
