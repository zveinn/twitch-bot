// Declare this file to be part of the main package so it can be compiled into
// an executable.
package main

// Import all Go packages required for this file.
import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Version is a constant that stores the Disgord version information.
const Version = "v0.0.0-alpha"

// Session is declared in the global space so it can be easily used
// throughout this program.
// In this use case, there is no error that would be returned.

func main() {
	var Session, _ = discordgo.New(os.Getenv("DISCORD_KEY"))
	// Declare any variables needed later.
	var err error

	// Print out a fancy logo!
	fmt.Printf(` 
	________  .__                               .___
	\______ \ |__| ______ ____   ___________  __| _/
	||    |  \|  |/  ___// ___\ /  _ \_  __ \/ __ | 
	||    '   \  |\___ \/ /_/  >  <_> )  | \/ /_/ | 
	||______  /__/____  >___  / \____/|__|  \____ | 
	\_______\/        \/_____/   %-16s\/`+"\n\n", Version)

	// Parse command line arguments
	flag.Parse()

	// Verify a Token was provided
	if Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Open a websocket connection to Discord
	err = Session.Open()
	defer Session.Close()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}
	for {
		time.Sleep(10 * time.Second)
		for _, guild := range Session.State.Guilds {

			// Get channels for this guild
			channels, _ := Session.GuildChannels(guild.ID)

			for _, c := range channels {
				// Check if channel is a guild text channel and not a voice or DM channel
				if c.Type != discordgo.ChannelTypeGuildText {
					continue
				}

				xxx, err := Session.ChannelMessages(c.ID, 10, "", latest, "")
				latest = xxx[0].ID
				for i, v := range xxx {
					log.Println(i, v)
				}
				log.Println(err)
				// Send text message
				// Session.ChannelMessageSend(
				// 	c.ID,
				// 	fmt.Sprintf("testmsg (sorry for spam). Channel name is %q", c.Name),
				// )
			}
		}
	}

	// // Wait for a CTRL-C
	// log.Printf(`Now running. Press CTRL-C to exit.`)
	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	// <-sc

	// // Clean up
	// Session.Close()

	// Exit Normally.
}

var latest string = "713220293308448828"
