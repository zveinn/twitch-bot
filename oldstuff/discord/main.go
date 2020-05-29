package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
)

func printMessage(session disgord.Session, evt *disgord.MessageCreate) {
	msg := evt.Message
	fmt.Println(msg.Author.String() + ": " + msg.Content) // Anders#7248{435358734985}: Hello there
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// see docs/examples/* for more information about configuration and use cases
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_KEY"),
		// BotToken: "dasdasdasdasdSDdas",
	})
	// connect, and stay connected until a system interrupt takes place
	defer client.StayConnectedUntilInterrupted(context.Background())

	log.Println("starting discord listener.")
	client.On(disgord.EvtMessageCreate, printMessage)
	msg, err := client.CreateMessage(context.Background(), disgord.Snowflake(710530877142335651), &disgord.CreateMessageParams{
		Content: "Test message from bot..",
	}, disgord.Flag(1<<4))
	if err != nil {
		log.Println(err, string(debug.Stack()))
	}
	log.Println(msg)

	// msg.Send(context.Background(), client, disgord.Flag(1<<4))
}
