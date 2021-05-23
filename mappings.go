package main

import (
	"github.com/andersfylling/disgord"
	"github.com/gempir/go-twitch-irc"
	"github.com/micmonay/keybd_event"
	"golang.org/x/net/websocket"
)

var keyboardCommandsToActions = make(map[string]int)
var socialCommandsToText = make(map[string]string)
var WSMAP = make(map[string]*websocket.Conn)
var oneMeter = 180
var TWITCHclient *twitch.Client
var DISCORDclient *disgord.Client
var KEYBONDING keybd_event.KeyBonding
var DISCORDCHANNEL uint64 = 713485972821639289
var TWITCHCHANNEL = "zzveinn"

// 500 ...
// 45 daagrees in the turn
// 2 meters foward
// 1 meter backwards
func LoadMaps() {
	// Media controls
	keyboardCommandsToActions["!prev"] = keybd_event.VK_MEDIA_PREV_TRACK
	keyboardCommandsToActions["!next"] = keybd_event.VK_MEDIA_NEXT_TRACK
	keyboardCommandsToActions["!play"] = keybd_event.VK_MEDIA_PLAY_PAUSE
	keyboardCommandsToActions["!volumeup"] = keybd_event.VK_VOLUME_UP
	keyboardCommandsToActions["!volumedown"] = keybd_event.VK_VOLUME_DOWN

	// generic controls
	keyboardCommandsToActions["!w"] = keybd_event.VK_W
	keyboardCommandsToActions["!a"] = keybd_event.VK_A
	keyboardCommandsToActions["!s"] = keybd_event.VK_S
	keyboardCommandsToActions["!d"] = keybd_event.VK_D

	// world of warcraft specific
	// keyboardCommandsToActions["!inventory"] = keybd_event.VK_B
	// keyboardCommandsToActions["!mount"] = keybd_event.VK_V
	// keyboardCommandsToActions["!jump"] = keybd_event.VK_SPACE
	// keyboardCommandsToActions["!dance"] = keybd_event.VK_K

	// Social stuff
	socialCommandsToText["!twitter"] = "https://www.twitter.com/zzveinn"
	socialCommandsToText["!youtube"] = "https://www.youtube.com/channel/UCW6eiMiVqYroPX1qiosAbnQ"
	socialCommandsToText["!discord"] = "https://discord.gg/r4wxkXd"
	socialCommandsToText["!hi"] = "Welcome to the stream, whats up!"

}
