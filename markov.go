package main

import (
	"log"
	"strings"

	"github.com/garyburd/redigo/redis"
)

type Markov struct {
	length int
}

func (m Markov) Store(text string, c redis.Conn) {
	text = strings.ToLower(text)
	splitted := strings.Split(text, " ")

	for k, v := range splitted {
		if k < len(splitted)-1 {
			c.Do("SADD", v, splitted[k+1])
		}

	}

}

func (m Markov) Generate(seed string, c redis.Conn) string {
	log.Printf("seed: %s\n", seed)

	seed = strings.ToLower(seed)
	splitted := strings.Split(seed, " ")

	key := string(splitted[0])

	s := []string{}

	if len(splitted) > 2 {
		for i := 1; i < m.length; i++ {
			s = append(s, key)

			next, _ := redis.String(c.Do("SRANDMEMBER", key))
			if next == "" {
				break
			}

			key = next
		}
	}
	text := strings.Join(s, " ")
	log.Printf("text: %s\n", text)
	return text
}
