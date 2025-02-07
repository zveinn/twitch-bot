package main

import (
	"bufio"
	"fmt"
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
		RenewTokens()
		time.Sleep(1 * time.Hour)
	}
}

func RenewTokens() error {
	file, err := os.Open("XR")
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	refreshtoken := scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Println(err)
		return err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	})
	if err != nil {
		log.Println("ERROR MAKIGN NEW HELIX CLIENT")
		log.Println(err)
		return err
	}
	resp, err := client.RefreshUserAccessToken(refreshtoken)
	if err != nil {
		log.Println("ERROR REFRESHING CREDENTIALS")
		log.Println(err)
		return err
	}
	os.Remove("X")
	keyFile, err := os.Create("X")
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("RESP:", resp.Data.AccessToken)
	if resp.Data.AccessToken == "" {
		os.Setenv("TWITCH_KEY", "oauth:"+refreshtoken)
		return nil
	}
	keyFile.WriteString(resp.Data.AccessToken)
	fmt.Println("TOKEN: ", "oauth:"+resp.Data.AccessToken)
	os.Setenv("TWITCH_KEY", "oauth:"+resp.Data.AccessToken)
	keyFile.Close()

	os.Remove("XR")
	refreshFile, err := os.Create("XR")
	if err != nil {
		log.Println(err)
		return err
	}
	refreshFile.WriteString(resp.Data.RefreshToken)
	refreshFile.Close()
	return nil
}

func MakeNewToken() {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURI:  "http://localhost:3000",
	})
	if err != nil {
		log.Println("ERROR MAKIGN NEW HELIX CLIENT")
		log.Println(err)
		return
	}
	// token := client.GetUserAccessToken()
	urlP := new(helix.AuthorizationURLParams)
	urlP.Scopes = append(urlP.Scopes, "channel:bot", "chat:edit", "chat:read", "user:bot", "user:read:chat", "user:write:chat", "whispers:read", "whispers:edit")
	urlP.ResponseType = "token"
	authUrl := client.GetAuthorizationURL(urlP)
	fmt.Println(authUrl)
	// resp, err := client.RequestUserAccessToken("ABCDDD")
	// fmt.Println(err)
	// fmt.Println(resp)
	// fmt.Println(resp.Data.AccessToken)
}
