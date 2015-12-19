package main

import (
	"log"
	"strings"

	"github.com/garyburd/redigo/redis"
)

type Markov struct {
	length int
}


func (m Markov) StoreUpdates(updates []Result, connection redis.Conn) {
	for _, update := range updates {
		if update.Message.Text != "" {
			m.Store(update.Message.Text, connection)
		}
	}
}

func (m Markov) Store(text string, c redis.Conn) {
	text = strings.ToLower(text) //todo, no lower case
	splitted := strings.Split(text, " ")

	for index, word := range splitted {
		if index < len(splitted)-1 {
			c.Do("SADD", word, splitted[index+1])
		}
	}

}

func (m Markov) Generate(seed string, connection redis.Conn) string {
	log.Printf("seed: %s\n", seed)

	seed = strings.ToLower(seed)
	splitted := strings.Split(seed, " ")

	key := string(splitted[0])

	s := []string{}

	if len(splitted) > 2 {
		for i := 1; i < m.length; i++ {
			s = append(s, key)

			next, _ := redis.String(connection.Do("SRANDMEMBER", key))
			if next == "" {
				break
			}

			key = next
		}
	}
	
	text := strings.Join(s, " ")
	text = text + "."
	log.Printf("Text: %s\n", text)

	return text
}
