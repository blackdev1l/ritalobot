package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const APIURL = "https://api.telegram.org/bot"

type Bot struct {
	Token      string
	ChatID     int
	Connection redis.Conn
}

//this should only return updates, it should not store them!!
func (bot Bot) GetUpdates() []Result {

	var jsonResp Response

	offset, _ := redis.String(bot.Connection.Do("GET", "update_id"))

	log.Printf("Getting Updates")
	resp, err := http.Get(APIURL + token + "/getUpdates?offset=" + offset)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("Parsing Updates")
	tempUpdates, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal(tempUpdates, &jsonResp)

	var updates = jsonResp.Result
	log.Printf("%v messages downloaded\n", len(updates))

	updateID := strconv.Itoa(updates[len(updates)-1].Update_id)
	bot.Connection.Do("SET", "update_id", updateID) //TODO: update id + 1, download 0 msgs

	return jsonResp.Result
}

func (bot Bot) Say(text string) int {
	chat := strconv.Itoa(chatID)
	resp, err := http.Get(APIURL + token + "/sendMessage?chat_id=" + chat + "&text=" + text)
	if err != nil {
		log.Println(err)
	}
	return resp.StatusCode
}

func (bot Bot) Connect(connection string, p int) {
	var err error

	port := ":" + strconv.Itoa(p)
	bot.Connection, err = redis.Dial(connection, port)
	if err != nil {
		fmt.Println("connection to redis failed")
		log.Fatal(err)
	}
	fmt.Printf("redis connection: %v | port is %v\n", connection, port)

}

func (bot Bot) Run() {
	timerUpdates := time.NewTicker(30 * time.Second)
	timerMessage := time.NewTicker(5 * time.Minute)

	markov := Markov{10}

	quit := make(chan struct{})

	var seed string

	for {
		select {
		case <-timerUpdates.C:
			var updates = bot.GetUpdates()

			markov.StoreUpdates(updates, bot.Connection)

			seed = updates[len(updates)-1].Message.Text
			fmt.Printf("Next Seed: %s", seed)
			break

		case <-timerMessage.C:
			text := markov.Generate(seed, bot.Connection)
			bot.Say(text)
			break

		case <-quit:
			timerMessage.Stop()
			timerUpdates.Stop()
			bot.Connection.Close()
			return
		}
	}
}
