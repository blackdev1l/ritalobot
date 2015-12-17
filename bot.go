package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const APIURL = "https://api.telegram.org/bot"

type Bot struct {
	Token  string
	ChatID int
	Seed   string
	C      redis.Conn
}

func (b Bot) GetUpdates(m Markov) {
	var esp Response

	offset, _ := redis.String(b.C.Do("GET", "update_id"))
	resp, err := http.Get(APIURL + b.Token + "/getUpdates?offset=" + offset)
	if err != nil {
		fmt.Println(err)
	}

	updates, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal(updates, &esp)

	for _, v := range esp.Result {
		if v.Message.Text != "" {

			m.Store(v.Message.Text, b.C)
		}

		updateID := strconv.Itoa(esp.Result[len(esp.Result)-1].Update_id)

		b.C.Do("SET", "update_id", updateID)

	}
	b.Seed = esp.Result[len(esp.Result)-1].Message.Text
}

func (b Bot) Say(text string) {
	_, err := http.Get(APIURL + b.Token + "/sendMessage?chat_id" + string(chatID) + "&text=" + text)
	if err != nil {
		fmt.Println(err)
	}
}

func (b Bot) Init() {

	var err error

	port := ":" + strconv.Itoa(port)
	b.C, err = redis.Dial(connection, port)
	if err != nil {
		fmt.Println("connection to redis failed")
		log.Fatal(err)
	}
	defer b.C.Close()

	b.Token = token
}

func (b Bot) Run() {
	mark := time.NewTicker(4 * time.Minute)
	s := time.NewTicker(5 * time.Minute)

	markov := Markov{10}

	quit := make(chan struct{})

	for {
		select {
		case <-mark.C:
			b.GetUpdates(markov)
			break
		case <-s.C:

			text := markov.Generate(b.Seed, b.C)
			b.Say(text)
			break

		case <-quit:
			mark.Stop()
			return
		}
	}
}
