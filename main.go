package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	oidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
)


func main () {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	clientID := os.Getenv("clientID")
	clientSecret := os.Getenv("clientSecret")

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/myrealm")
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8081/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}


	state := string(rand.Intn(100))

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Query().Get("state") != state {
			http.Error(writer, "State invalido", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, request.URL.Query().Get("code"))
		if err != nil {
			http.Error(writer, "Falha ao trocar o token", http.StatusInternalServerError)
			return
		}

		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(writer, "Falha ao gerar o IDToken", http.StatusInternalServerError)
			return
		}

		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
		if err != nil {
			http.Error(writer, "Falha ao pegar UserInfo", http.StatusInternalServerError)
			return
		}

		type User struct {
			AcessToken *oauth2.Token
			IDToken    string
			UserInfo   *oidc.UserInfo
		}

		resp := User{
			token,
			idToken,
			userInfo,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Write(data)



	})

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {

		rawAccessToken := request.Header.Get("Authorization")
		output := fmt.Sprint(rawAccessToken)

		var validRole = regexp.MustCompile(`bV9hY2Nlc3MiOnsicm9sZXMiOlsiYWRtaW4iXX0s`)

		validRole.MatchString(output)

		r,_ := base64.StdEncoding.DecodeString(validRole.FindString(output))

		var validStr = regexp.MustCompile(`{(.*?)\}`)
		role := validStr.FindString(string(r))

		var result map[string]interface{}
		json.Unmarshal([]byte(role), &result)

		json_obj := result["roles"].([]interface{})
		fmt.Printf("%s",json_obj)

		for _, value := range json_obj {

			fmt.Println(value.(string))

			if value.(string) != "admin"{
				http.Error(writer, "Forbidden", http.StatusForbidden)
				return
			}

			if  value.(string) == "admin"{
				writer.Write([]byte("hello world"))
				return
			}

		}



	})

	http.HandleFunc("/lock", func(writer http.ResponseWriter, request *http.Request) {

		rawAccessToken := request.Header.Get("Authorization")
		output := fmt.Sprint(rawAccessToken)


		var validRole = regexp.MustCompile(`bV9hY2Nlc3MiOnsicm9sZXMiOlsiYWRtaW4iXX0s`)

		validRole.MatchString(output)
		fmt.Println(validRole.FindString(output))


		r,_ := base64.StdEncoding.DecodeString(validRole.FindString(output))
		sr := string(r)

		var validStr = regexp.MustCompile(`{(.*?)\}`)
		role := validStr.FindString(sr)

		var result map[string]interface{}
		json.Unmarshal([]byte(role), &result)

		json_obj := result["roles"].([]interface{})

		for _, value := range json_obj {

			fmt.Println(value.(string))

			if value.(string) != "test"{
				http.Error(writer, "Forbidden", http.StatusForbidden)
				return
			}

			if  value.(string) == "test"{
				writer.Write([]byte("lock"))
				return
			}

		}


	})


	log.Fatal(http.ListenAndServe(":8081", nil))



}


