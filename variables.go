package main

import tirc "github.com/gempir/go-twitch-irc/v4"

// FILES
// X == oauth token
// XR == oauth refresh token

// REQUIRED .ENV VARIABLES
// TWITCH_KEY --- this variable is set automatically
// CLIENT_ID  --- client ID for twitch API
// CLIENT_SECRET --- client secret for twitch API
// DB --- the database connection string to your mongoDB

var (
	USERNAME       = "keyb1nd_"
	BROADCASTER_ID = "704389637"
)

func InitMP3Map() {
	MP3Map["e7364edf-c725-45e5-9938-3fbd7659fd07"] = "herewegoagain"
	MP3Map["ae16a19b-5bea-4407-838d-c2895d41db6c"] = "ohmygod"
	MP3Map["c63ecae5-2183-4b39-a5a2-6be8150732ea"] = "thisistheway"
	MP3Map["d60c042f-3540-4fc3-aa8a-e42617815880"] = "developers"
	MP3Map["cd14e9cd-73fc-4dec-9c60-6e267082351f"] = "aya-short"
	MP3Map["ebe70cf1-340e-4753-872f-281aa49f8505"] = "uwu-long"

	MP3Map["14537722-6db9-4e89-bc3e-89c7e7e19e91"] = "uwu"
	MP3Map["1ef58611-e4b4-4f9c-a96e-e8a333b1e0e7"] = "araara"
	MP3Map["aadbb1ef-c995-4ca8-b720-8dc7758389ce"] = "onichan"
	MP3Map["0878e53e-67cd-47d1-89a1-91cc241d9412"] = "gnupluslinux"

	MP3Map["2a174cdd-a444-434e-af31-9e6a598944de"] = "come-after-you"
	MP3Map["5a78d388-6757-422b-a348-9ce983f34cb3"] = "hey-listen"
	MP3Map["1594a455-4a84-4a8b-a562-ac830c423d81"] = "excellent"
	MP3Map["b714ac5a-e6a4-4cf5-9ea0-385c67a168b8"] = "blabla"
}

func InitTwitchClient() {
	TWITCH_CLIENT.Name = USERNAME
	TWITCH_CLIENT.ChannelMap = make(map[string]*IRC_CHANNEL)
	TWITCH_CLIENT.ChannelMap[USERNAME] = new(IRC_CHANNEL)
	TWITCH_CLIENT.ChannelMap[USERNAME].BroadCasterID = BROADCASTER_ID
	TWITCH_CLIENT.ChannelMap[USERNAME].Name = USERNAME
	TWITCH_CLIENT.ChannelMap[USERNAME].Type = 1
}

func InitCommands() {
	// TextCommands["!monero"] = "43V6N2BpjvMYUthyqLioafZ2MQQviWEhvVTpp3hHc6LB48WYE8SsjrJKyyYzR3AYu2HkSXu8xsJhr7wdLsgSc8mGDDTkCrn"
	TextCommands["!nvim"] = "https://github.com/zveinn/dotfiles"
	TextCommands["!dotfiles"] = "https://github.com/zveinn/dotfiles"
	TextCommands["!x"] = "https://x.com/keyb1nd"
	TextCommands["!github"] = "https://github.com/zveinn"
	TextCommands["!linkedin"] = "https://www.linkedin.com/in/keyb1nd/"
	TextCommands["!discord"] = "https://discord.com/invite/wJ5m3Y6ezq"
	TextCommands["!keyboard"] = "wooting.io"
	TextCommands["!os"] = "All of them."
	TextCommands["!terminal"] = "wezterm + tmux"
	TextCommands["!editor"] = "nvim"
	TextCommands["!youtube"] = "https://www.youtube.com/@keyb1nd?sub_confirmation=1"
	TextCommands["!lurk"] = "ABSOLUTELY NOT ... LURKING IS NOT ALLOWED IN HERE"
	TextCommands["!signal"] = "https://signal.group/#CjQKILCHWDqtfKErs-6yV8i0kQHhScDTL4wQ2mW7JYoQoBLsEhC7R4AqmLxdxwdRa0fWK1tD"

	// VPN RELATED
	TextCommands["!freetrial"] = "All new accounts get 24 hours free trial > https://www.tunnels.is"
	TextCommands["!vpn"] = "Tunnels.is >> Advanced Networking Utility >> 24/h Free Trial >> https://tunnels.is"
	TextCommands["!vpndiscord"] = "Tunnels.is DISCORD: https://discord.com/invite/7Ts3PCnCd9"

	TextCommands["!cmd"] += "!tts !roll !points !quote !top10"
	TextCommands["!commands"] += "!tts !roll !points !quote !top10"

	for i := range TextCommands {
		TextCommands["!cmd"] += " " + i
		TextCommands["!commands"] += " " + i
	}

}

func CheckCustomReward(U *User, msg tirc.PrivateMessage) (success bool) {
	if msg.CustomRewardID == "8444968a-be3c-4d89-b6e7-dbbdedf64e1f" {
		go PlaceEventOnSoundQueue("tts", msg.Message)
		return true
	}
	if msg.CustomRewardID == "323be4d7-63e6-4f2d-ad99-246f19c9ebd7" {
		_ = IncrementUserPoints(U, 100)
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "Redeemed 100 points!", "")
		return true
	} else if msg.CustomRewardID == "601576ec-b3ad-4f2d-8bba-a8b79f2f7e14" {
		_ = IncrementUserPoints(U, 500)
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "Redeemed 500 points!", "")
		return true
	} else if msg.CustomRewardID == "a8b676f0-0dc3-441d-a2dc-7cb0de3499ee" {
		_ = IncrementUserPoints(U, 1000)
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "Redeemed 1000 points!", "")
		return true
	} else if msg.CustomRewardID == "a1b6ad63-a3da-492d-b37f-ad068997bd70" {
		_ = IncrementUserPoints(U, 5000)
		TWITCH_CLIENT.ReplyToUser(msg.User.DisplayName, "Redeemed 5000 points!", "")
		return true
	}
	return
}
