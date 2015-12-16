package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	chatID     int
	port       int
	token      string
	connection string
	configPath string
)

type Config struct {
	Token      string `yaml:"token"`
	ChatID     int    `yaml:"chatID"`
	Connection string `yaml:"connection"`
	Port       int    `yaml:"port"`
}

func printLogo() {
	file, _ := ioutil.ReadFile("logo")
	fmt.Println(string(file))
}

func readConfig() {

	filename, _ := filepath.Abs(configPath)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("no config file found, program will stop")
	}

	var c Config

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Println("error while parsing config.yaml, closing")
		log.Fatalln(err)
	}

	token = c.Token
	chatID = c.ChatID
	connection = c.Connection
	port = c.Port
}

func main() {

	printLogo()

	flag.StringVar(&token, "token", "", "authentication token for the telegram bot")
	flag.IntVar(&chatID, "id", 0, "Chat id of the group chat")
	flag.StringVar(&connection, "conn", "tcp", "type of connection and/or ip of redis database")
	flag.IntVar(&port, "p", 6379, "port number of redis database")
	flag.StringVar(&configPath, "c", "./config.yml", "path for ritalobot config")

	flag.Parse()

	readConfig()

	if token == "" || chatID == 0 {
		log.Fatalln("authentication token or chat id not valid, use config or flags to pass it")
	}

	bot := Bot{}
	bot.Init()

	bot.Run()

}
