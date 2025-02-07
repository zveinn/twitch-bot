package main

import (
	"fmt"
	"os/exec"
)

func main() {
	out, err := exec.Command("./ffmpeg.exe", "-f", "dshow", "-i", `audio=Microphone (Yeti Classic)`, "output.mp3").CombinedOutput()
	fmt.Println(string(out))
	fmt.Println(err)
}
