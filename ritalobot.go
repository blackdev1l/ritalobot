package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readToken() string {
	var token string

	file, err := os.Open("token")
	defer file.Close()
	if err != nil {
		log.Fatal("Make sure to have your token in a file called 'token'")
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		token += scanner.Text()
	}

	fmt.Printf("token is %s\n", token)

	return token

}
func main() {
	token := readToken()
	bot := Bot{}
	bot.Init(token)

	bot.Run()

}
