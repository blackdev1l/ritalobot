package main

import (
	"./ui"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	port       int
	token      string
	connection string
	configPath string
	chance     int
)

type Config struct {
	Token      string `yaml:"token"`
	Chance     int    `yaml: "chance"`
	Connection string `yaml:"connection"`
	Port       int    `yaml:"port"`
}

func printLogo() {
	file, _ := ioutil.ReadFile("logo")
	fmt.Println(string(file))
}

func readConfig(configPath string) int {

	filename, _ := filepath.Abs(configPath)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("no config file found")
		log.Println(err)
		return 1
	}

	var c Config

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Println("error while parsing config.yaml, closing")
		log.Fatalln(err)
	}

	token = c.Token
	chance = c.Chance
	connection = c.Connection
	port = c.Port
	return 0
}

func main() {
	botChan := make(chan int)

	//printLogo()

	flag.StringVar(&token, "token", "", "authentication token for the telegram bot")
	flag.StringVar(&connection, "conn", "tcp", "type of connection and/or ip of redis database")
	flag.IntVar(&port, "p", 6379, "port number of redis database")
	flag.StringVar(&configPath, "c", "./config.yml", "path for ritalobot config")
	flag.IntVar(&chance, "chance", 10, "chance to say something after a message")

	flag.Parse()

	readConfig(configPath)

	if token == "" {
		log.Fatalln("authentication token not valid, use config or flags to pass it")
	}

	go ui.Show(botChan)
	bot := Bot{}
	bot.Listen(botChan)

}
