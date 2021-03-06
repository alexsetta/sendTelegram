package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	token  string
	chatID int64
}

// The main funtion starts our server on port 3000
func main() {
	if len(os.Args) != 3 {
		log.Fatal("parâmetros inválidos")
	}
	base := os.Args[1]
	file := os.Args[2]

	b, err := ioutil.ReadFile(base + "/sendTelegram.cfg")
	if err != nil {
		log.Fatal(err)
	}

	lido := strings.Split(string(b), "|")
	if len(lido) != 2 {
		log.Fatal("arquivo de configuração inválido")
	}

	i, err := strconv.ParseInt(lido[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	c := Config{lido[0], i}

	b, err = ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	if err := sendMessage(c, string(b)); err != nil {
		log.Fatal(err)
	}
}

// Create a struct to conform to the JSON body of the send message request
// https://core.telegram.org/bots/api#sendmessage
type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func sendMessage(c Config, message string) error {
	// Create the request body struct
	reqBody := &sendMessageReqBody{
		ChatID: c.chatID,
		Text:   message,
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot"+c.token+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status:" + res.Status)
	}

	return nil
}
