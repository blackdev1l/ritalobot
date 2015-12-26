package main

import (
	"os"
	"strconv"
	"testing"
)

var bot Bot

func TestSay(t *testing.T) {
	test := "BEEP BOOP THIS IS A TEST"
	if bot.Say(test) != 200 {
		t.Fail()
	}
}

func TestMain(m *testing.M) {

	if readConfig("./config.yml") != 0 {
		tmp, _ := strconv.Atoi(os.Getenv("chat"))
		chatID = tmp
		token = os.Getenv("token")
	}

	os.Exit(m.Run())
}
