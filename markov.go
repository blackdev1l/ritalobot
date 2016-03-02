package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"regexp"
	"strings"
)

type Markov struct {
	length int
}

func (m Markov) StoreUpdates(updates []Result, connection redis.Conn) {
	for _, update := range updates {
		message := update.Message.Text
		if message != "" && !strings.HasPrefix(message, "/") {
			m.Store(update.Message.Text, connection)
		}
	}
}

func (m Markov) Store(text string, c redis.Conn) {
	splitted := strings.Split(text, " ")

	for index, word := range splitted {
		if index < len(splitted)-1 {
			c.Do("SADD", word, splitted[index+1])
		}
	}
}

func (m Markov) Generate(seed string, connection redis.Conn) string {
	log.Printf("seed: %s\n", seed)

	s := []string{}

	s = append(s, seed)

	splitted := strings.Split(seed, " ")

	// Start the chain with the last word of the seed
	key := splitted[len(splitted)-1]

	for i := 1; i < m.length; i++ {

		next, _ := redis.String(connection.Do("SRANDMEMBER", key))
		s = append(s, next)

		matched, _ := regexp.MatchString(".*[\\.;!?¿¡]$", next)
		if next == "" || matched {
			break
		}

		key = next
	}

	text := strings.Join(s, " ")
	log.Printf("Text: %s\n", text)
	return text
}
