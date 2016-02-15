package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

func sendCommand(method, token string, params url.Values) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s?%s",
		token, method, params.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	resp.Close = true
	defer resp.Body.Close()
	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return json, nil
}

const APIURL = "https://api.telegram.org/bot"

type Bot struct {
	Token      string
	ChatID     int
	Connection redis.Conn
}

//this should only return updates, it should not store them!!
func (bot Bot) GetUpdates() []Result {

	timeout := 35 * time.Second
	var jsonResp Response

	client := http.Client{
		Timeout: timeout,
	}

	offset, _ := redis.String(bot.Connection.Do("GET", "update_id"))

	log.Printf("Getting Updates")
	resp, err := client.Get(APIURL + token + "/getUpdates?offset=" + offset + "&timeout=" + strconv.Itoa(30))
	if err != nil {
		fmt.Println(err)
	}

	tempUpdates, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal(tempUpdates, &jsonResp)

	log.Printf("Parsing Updates")
	var updates = jsonResp.Result
	if len(updates) != 0 {
		log.Printf("%v messages downloaded\n", len(updates))

		updateID := updates[len(updates)-1].Update_id + 1
		bot.Connection.Do("SET", "update_id", updateID) //TODO: update id + 1, download 0 msgs

		return jsonResp.Result

	}
	return nil
}

func (bot Bot) Say(text string) int {
	chat := strconv.Itoa(chatID)
	params := url.Values{}

	params.Set("chat_id", chat)
	params.Set("text", text)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s?%s",
		token, "sendMessage", params.Encode())
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	return resp.StatusCode
}

func (bot Bot) Listen() {
	var err error
	var seed string
	markov := Markov{10}
	rand.Seed(time.Now().UnixNano())

	tmp := ":" + strconv.Itoa(port)
	bot.Connection, err = redis.Dial(connection, tmp)
	if err != nil {
		fmt.Println("connection to redis failed")
		log.Fatal(err)
	}
	fmt.Printf("redis connection: %v | port is %v\n", connection, port)

	for {
		updates := bot.GetUpdates()
		if updates != nil {
			markov.StoreUpdates(updates, bot.Connection)
			if rand.Intn(100) <= 10 {
				seed = updates[len(updates)-1].Message.Text
				fmt.Printf("Next Seed: %s", seed)
				text := markov.Generate(seed, bot.Connection)
				bot.Say(text)

			}
		}
	}
}
