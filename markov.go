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
		if update.Message.Text != "" {
			m.Store(update.Message.Text, connection)
		}
	}
}

func (m Markov) Store(text string, c redis.Conn) {
	splitted := strings.Split(text, " ")

	// if the first word has '/' character skip the whole string
	if !strings.ContainsAny(splitted[0], "/") {
		for index, word := range splitted {
			if index < len(splitted)-1 {
				c.Do("SADD", word, splitted[index+1])
			}
		}
	}
}

func (m Markov) Generate(seed string, connection redis.Conn) string {
	log.Printf("seed: %s\n", seed)

	seed = strings.ToLower(seed)
	splitted := strings.Split(seed, " ")

	key := string(splitted[0])

	s := []string{}

	for i := 1; i < m.length; i++ {
		s = append(s, key)

		next, _ := redis.String(connection.Do("SRANDMEMBER", key))
		key = next

		matched, _ := regexp.MatchString(".*[\\.;!?¿¡]$", next)
		if next == "" || matched {
			break
		}
	}

	text := strings.Join(s, " ")
	log.Printf("Text: %s\n", text)
	return text
}
