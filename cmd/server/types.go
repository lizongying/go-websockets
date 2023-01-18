package main

type Message struct {
	Id   int    `json:"id"`
	From string `json:"from"`
	Via  string `json:"via"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
}
