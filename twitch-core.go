package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"

	"github.com/nicklaw5/helix"
)

func SetTwitchKeyEnvVariable() {
	tk, err := os.Open("X")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	tkbyte, err := io.ReadAll(tk)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Setenv("TWITCH_KEY", "oauth:"+string(tkbyte))
}

func RenewTokensLoop() {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r)
		}
		monitor <- 7
	}()

	for {
		time.Sleep(1 * time.Hour)
		RenewTokens()
	}
}

func RenewTokens() {
	file, err := os.Open("XR")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	refreshtoken := scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Println(err)
		return
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	})
	if err != nil {
		log.Println("ERROR MAKIGN NEW HELIX CLIENT")
		log.Println(err)
		return
	}
	resp, err := client.RefreshUserAccessToken(refreshtoken)
	if err != nil {
		log.Println("ERROR REFRESHING CREDENTIALS")
		log.Println(err)
		return
	}
	os.Remove("X")
	keyFile, err := os.Create("X")
	if err != nil {
		log.Println(err)
		return
	}
	keyFile.WriteString(resp.Data.AccessToken)
	os.Setenv("TWITCH_KEY", "oauth:"+resp.Data.AccessToken)
	keyFile.Close()

	os.Remove("XR")
	refreshFile, err := os.Create("XR")
	if err != nil {
		log.Println(err)
		return
	}
	refreshFile.WriteString(resp.Data.RefreshToken)
	refreshFile.Close()
}
