package main

import (
	"fmt"
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
	fmt.Printf("seed is %s\n", seed)
	var text string

	seed = strings.ToLower(seed)
	splitted := strings.Split(seed, " ")

	key := string(splitted[0])

	if len(splitted) > 2 {
		for i := 1; i < m.length; i++ {
			text = text + " " + key

			next, _ := redis.String(c.Do("SRANDMEMBER", key))
			if next == "" {
				return text
			}

			key = next
		}
	}

	return text
}
