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
	"regexp"
	"github.com/garyburd/redigo/redis"
)

func sendCommand(method, token string, params url.Values) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s?%s",
		token, method, params.Encode())

	timeout := 35 * time.Second

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
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

type Bot struct {
	Token      string
	Connection redis.Conn
	Chance     int
}

func (bot Bot) GetUpdates() []Result {
	offset, _ := redis.String(bot.Connection.Do("GET", "update_id"))

	params := url.Values{}
	params.Set("offset", offset)
	params.Set("timeout", strconv.Itoa(30))

	resp, err := sendCommand("getUpdates", token, params)
	if err != nil {
		log.Println(err)
	}

	var updatesRecieved Response
	json.Unmarshal(resp, &updatesRecieved)

	if !updatesRecieved.Ok {
		err = fmt.Errorf("chobot: %s\n", updatesRecieved.Description)
		return nil
	}

	var updates = updatesRecieved.Result
	if len(updates) != 0 {
		log.Printf("%v messages downloaded\n", len(updates))

		updateID := updates[len(updates)-1].Update_id + 1
		bot.Connection.Do("SET", "update_id", updateID)

		return updates

	}
	return nil
}

func (bot Bot) Say(text string, chat int) (bool, error) {

	var responseRecieved struct {
		Ok          bool
		Description string
	}

	params := url.Values{}

	params.Set("chat_id", strconv.Itoa(chat))
	params.Set("text", text)
	resp, err := sendCommand("sendMessage", token, params)

	err = json.Unmarshal(resp, &responseRecieved)
	if err != nil {
		return false, err
	}

	if !responseRecieved.Ok {
		return false, fmt.Errorf("chobot: %s\n", responseRecieved.Description)
	}

	return responseRecieved.Ok, nil
}

func (bot Bot) Listen() {
	var err error
	var seed string
	markov := Markov{10}
	rand.Seed(time.Now().UnixNano())
	bot.Chance = chance

	tmp := ":" + strconv.Itoa(port)
	bot.Connection, err = redis.Dial(connection, tmp)
	if err != nil {
		fmt.Println("connection to redis failed")
		log.Fatal(err)
	}
	fmt.Printf("redis connection: %v | port is %v\n", connection, port)
	fmt.Printf("chance rate %v%!\n", bot.Chance)

	for {
		updates := bot.GetUpdates()
		if updates != nil {

			hasSpoken :=false
			for i := 0; i<len(updates); i++ {
				update := updates[i]
				re := regexp.MustCompile("\\/chobot\\s+(.+?)\\b")
				match := re.FindAllStringSubmatch(update.Message.Text, -1)
				if len(match)>0{
					text := markov.Generate(match[0][1], bot.Connection)
					bot.Say(text, update.Message.Chat.Id)
					hasSpoken = true
				}
				updates = append(updates[:i], updates[i+1:]...)	
			}

			if len(updates)>0 {
				markov.StoreUpdates(updates, bot.Connection)
				if rand.Intn(100) <= bot.Chance && !hasSpoken {
					seed = updates[len(updates)-1].Message.Text
					chat := updates[len(updates)-1].Message.Chat.Id
					fmt.Printf("Next Seed: %s", seed)
					text := markov.Generate(seed, bot.Connection)
					bot.Say(text, chat)
				}
			}
			
		}
	}
}
