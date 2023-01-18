package main

import (
	"encoding/json"
	"flag"
	"golang.org/x/net/websocket"
	"log"
	"sync"
)

type Client struct {
	id     string
	url    string
	origin string
	lastId int
	ws     *websocket.Conn
	chs    sync.Map
}

func (c *Client) register() (err error) {
	c.ws, err = websocket.Dial(c.url, "", c.origin)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (c *Client) ping() (err error) {
	m := Message{
		From: c.id,
	}
	s, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = c.ws.Write(s)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (c *Client) send(via string, to string, msg string) (id int, err error) {
	c.lastId++
	id = c.lastId
	m := Message{
		Id:   id,
		From: c.id,
		Via:  via,
		To:   to,
		Msg:  msg,
	}
	s, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = c.ws.Write(s)
	if err != nil {
		log.Println(err)
		return
	}
	ch := make(chan Message, 1)
	c.chs.Store(id, ch)
	log.Printf("send: %+v\n", m)
	return
}

func (c *Client) keep() {
	for {
		var msg = make([]byte, 1024)
		var n int
		n, err := c.ws.Read(msg)
		if err != nil {
			if err.Error() == "EOF" {
				continue
			}
			log.Println(err)
			break
		}
		var m Message
		err = json.Unmarshal(msg[:n], &m)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("received: %+v\n", m)

		if m.Via == "" && m.From == m.To {
			value, ok := c.chs.Load(m.Id)
			if ok {
				value.(chan Message) <- m
			}
		}

		if m.Via != "" {
			m.Via = ""
			m.Msg = "ok"
			s, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
				break
			}
			_, err = c.ws.Write(s)
			if err != nil {
				log.Println(err)
				break
			}
			log.Printf("send: %+v\n", m)
		}
	}
}

func NewClient(url string, origin string, id string) (client *Client, err error) {
	client = &Client{
		id:     id,
		url:    url,
		origin: origin,
	}
	err = client.register()
	if err != nil {
		return
	}
	err = client.ping()
	if err != nil {
		return
	}
	go client.keep()
	return
}

func main() {
	urlPtr := flag.String("url", "ws://127.0.0.1:1234/echo", "--url=ws://127.0.0.1:1234/echo")
	originPtr := flag.String("origin", "http://127.0.0.1/", "--origin=http://127.0.0.1/")
	fromPtr := flag.String("from", "", "--from=client1/client2")
	viaPtr := flag.String("via", "", "--via=client1/client2")
	toPtr := flag.String("to", "", "--to=client1/client2")
	msgPtr := flag.String("msg", "", "--msg=hi")
	waitPtr := flag.Bool("wait", false, "--wait")
	flag.Parse()
	log.Printf("url: %s, origin: %s, from: %s\n", *urlPtr, *originPtr, *fromPtr)

	client, err := NewClient(*urlPtr, *originPtr, *fromPtr)
	if err != nil {
		log.Panicln(err)
	}
	if *toPtr != "" {
		if *waitPtr {
			id, err := client.send(*viaPtr, *toPtr, *msgPtr)
			if err != nil {
				log.Panicln(err)
			}
			value, ok := client.chs.Load(id)
			if ok {
				res := <-value.(chan Message)
				log.Printf("res: %+v", res)
				client.chs.Delete(id)
			}
		} else {
			_, err = client.send(*viaPtr, *toPtr, *msgPtr)
			if err != nil {
				log.Panicln(err)
			}
		}
	}
	select {}
}
