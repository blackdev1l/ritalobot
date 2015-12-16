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

const API_URL = "https://api.telegram.org/bot"

type Bot struct {
	Token   string
	Chat_id int
	Seed    string
	C       redis.Conn
}

func (b Bot) GetUpdates(m Markov) {

	var esp Response

	offset, _ := redis.String(b.C.Do("GET", "update_id"))
	resp, err := http.Get(API_URL + b.Token + "/getUpdates?offset=" + offset)
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

		update_id := strconv.Itoa(esp.Result[len(esp.Result)-1].Update_id)

		b.C.Do("SET", "update_id", update_id)

	}
	b.Seed = esp.Result[len(esp.Result)-1].Message.Text
}

func (b Bot) Say(text string) {
	_, err := http.Get(API_URL + b.Token + "/sendMessage?chat_id=-15689316&text=" + text)
	if err != nil {
		fmt.Println(err)
	}
}

func (b Bot) Init(token string) {

	var err error

	b.C, err = redis.Dial("tcp", ":6379")
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
