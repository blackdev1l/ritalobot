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
	C      redis.Conn
}

func (b Bot) GetUpdates(m Markov) string {
	var esp Response

	offset, _ := redis.String(b.C.Do("GET", "update_id"))
	resp, err := http.Get(APIURL + token + "/getUpdates?offset=" + offset)
	if err != nil {
		fmt.Println(err)
	}

	updates, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal(updates, &esp)

	log.Printf("%v messages downloaded\n", len(esp.Result))
	for _, v := range esp.Result {
		if v.Message.Text != "" {
			m.Store(v.Message.Text, b.C)
		}

		updateID := strconv.Itoa(esp.Result[len(esp.Result)-1].Update_id)

		b.C.Do("SET", "update_id", updateID)

	}
	return esp.Result[len(esp.Result)-1].Message.Text
}

func (b Bot) Say(text string) {
	chat := strconv.Itoa(chatID)
	fmt.Println(APIURL + token + "/sendMessage?chat_id=" + chat + "&text=" + text)
	_, err := http.Get(APIURL + token + "/sendMessage?chat_id" + chat + "&text=" + text)
	if err != nil {
		log.Println(err)
	}
	log.Println("sent")
}

func (b Bot) Run() {
	var err error

	port := ":" + strconv.Itoa(port)
	b.C, err = redis.Dial(connection, port)
	if err != nil {
		fmt.Println("connection to redis failed")
		log.Fatal(err)
	}
	defer b.C.Close()

	fmt.Printf("redis connection: %v | port is %v\n", connection, port)

	mark := time.NewTicker(5 * time.Minute)

	markov := Markov{10}

	quit := make(chan struct{})

	for {
		select {
		case <-mark.C:
			seed := b.GetUpdates(markov)
			text := markov.Generate(seed, b.C)
			b.Say(text)
			break

		case <-quit:
			mark.Stop()
			return
		}
	}
}
