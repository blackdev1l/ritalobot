package main

type Response struct {
	Ok          bool     `json:"ok"`
	Result      []Result `json:"result"`
	Description string
}

type Result struct {
	Update_id int     `json:"update_id"`
	Message   Message `json:"message"`
}

type Message struct {
	Message_id int    `json:"message_id"`
	From       From   `json:"from"`
	Text       string `json:"text"`
}

type From struct {
	Id       int
	Username string
}
