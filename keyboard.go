//go:build windows

package main

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/zveinn/twitch-bot/mic"
)

const (
	// Add more keys as constants here if needed
	VK_LBUTTON  = 0x01 // Left mouse button
	VK_RBUTTON  = 0x02 // Right mouse button
	VK_CANCEL   = 0x03 // Control-break processing
	VK_MBUTTON  = 0x04 // Middle mouse button (three-button mouse)
	VK_XBUTTON1 = 0x05 // X1 mouse button
	VK_XBUTTON2 = 0x06 // X2 mouse button
	VK_BACK     = 0x08 // BACKSPACE key
	VK_TAB      = 0x09 // TAB key
	VK_CLEAR    = 0x0C // CLEAR key
	VK_RETURN   = 0x0D // ENTER key
	VK_SHIFT    = 0x10 // SHIFT key
	VK_CONTROL  = 0x11 // CTRL key
	VK_MENU     = 0x12 // ALT key
	VK_PAUSE    = 0x13 // PAUSE key
	VK_CAPITAL  = 0x14 // CAPS LOCK key
	VK_ESCAPE   = 0x1B // ESC key
	VK_SPACE    = 0x20 // SPACEBAR
	VK_PRIOR    = 0x21 // PAGE UP key
	VK_NEXT     = 0x22 // PAGE DOWN key
	VK_END      = 0x23 // END key
	VK_HOME     = 0x24 // HOME key
	VK_LEFT     = 0x25 // LEFT ARROW key
	VK_UP       = 0x26 // UP ARROW key
	VK_RIGHT    = 0x27 // RIGHT ARROW key
	VK_DOWN     = 0x28 // DOWN ARROW key
	VK_SELECT   = 0x29 // SELECT key
	VK_PRINT    = 0x2A // PRINT key
	VK_EXECUTE  = 0x2B // EXECUTE key
	VK_SNAPSHOT = 0x2C // PRINT SCREEN key
	VK_INSERT   = 0x2D // INS key
	VK_DELETE   = 0x2E // DEL key
	VK_HELP     = 0x2F // HELP key

	// Keypad keys
	VK_NUMLOCK   = 0x90 // NUM LOCK key
	VK_SCROLL    = 0x91 // SCROLL LOCK key
	VK_NUMPAD0   = 0x60 // Numeric keypad 0 key
	VK_NUMPAD1   = 0x61 // Numeric keypad 1 key
	VK_NUMPAD2   = 0x62 // Numeric keypad 2 key
	VK_NUMPAD3   = 0x63 // Numeric keypad 3 key
	VK_NUMPAD4   = 0x64 // Numeric keypad 4 key
	VK_NUMPAD5   = 0x65 // Numeric keypad 5 key
	VK_NUMPAD6   = 0x66 // Numeric keypad 6 key
	VK_NUMPAD7   = 0x67 // Numeric keypad 7 key
	VK_NUMPAD8   = 0x68 // Numeric keypad 8 key
	VK_NUMPAD9   = 0x69 // Numeric keypad 9 key
	VK_MULTIPLY  = 0x6A // Multiply key
	VK_ADD       = 0x6B // Add key
	VK_SEPARATOR = 0x6C // Separator key
	VK_SUBTRACT  = 0x6D // Subtract key
	VK_DECIMAL   = 0x6E // Decimal key
	VK_DIVIDE    = 0x6F // Divide key

	// Function keys
	VK_F1  = 0x70 // F1 key
	VK_F2  = 0x71 // F2 key
	VK_F3  = 0x72 // F3 key
	VK_F4  = 0x73 // F4 key
	VK_F5  = 0x74 // F5 key
	VK_F6  = 0x75 // F6 key
	VK_F7  = 0x76 // F7 key
	VK_F8  = 0x77 // F8 key
	VK_F9  = 0x78 // F9 key
	VK_F10 = 0x79 // F10 key
	VK_F11 = 0x7A // F11 key
	VK_F12 = 0x7B // F12 key

	// Other keys
	VK_LWIN  = 0x5B // Left Windows key (Natural keyboard)
	VK_RWIN  = 0x5C // Right Windows key (Natural keyboard)
	VK_APPS  = 0x5D // Applications key (Natural keyboard)
	VK_SLEEP = 0x5F // Computer Sleep key
)

func captureKeys() {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r)
		}
		monitor <- 13
	}()

	user32 := syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState := user32.NewProc("GetAsyncKeyState")

	fmt.Println("Listening for key presses... Press Ctrl+C to exit.")

	for {
		for _, vkCode := range []int{VK_XBUTTON2} {
			r1, _, _ := getAsyncKeyState.Call(uintptr(vkCode))
			if r1&0x8000 != 0 { // Check if the key is down
				fmt.Printf("Key %#x is down\n", vkCode)
				x := mic.TalkToEve()
				fmt.Println("POST TRANSCRIBE: ", x)
				PlaceBotEventInQueue("eve", x, x)
			}
		}
		time.Sleep(time.Millisecond * 50) // To prevent excessive CPU usage
	}
}
