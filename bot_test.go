package main

import (
	"os"
	"testing"
)

var bot Bot
var chat string

func TestMain(m *testing.M) {

	if readConfig("./config.yml") != 0 {
		token = os.Getenv("token")
	}

	os.Exit(m.Run())
}
