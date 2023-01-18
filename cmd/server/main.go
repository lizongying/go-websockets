package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

var users sync.Map

func handle(ws *websocket.Conn) {
	for {
		var msg = make([]byte, 1024)
		var n int
		n, err := ws.Read(msg)
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

		if m.Id == 0 {
			users.Store(m.From, ws)
		}

		if m.Msg != "" && m.To != "" {
			if m.Via != "" {
				value, ok := users.Load(m.Via)
				if ok {
					conn := value.(*websocket.Conn)
					m.Via = ""
					s, err := json.Marshal(m)
					if err != nil {
						log.Println(err)
						continue
					}
					log.Printf("send: %+v\n", m)
					_, err = conn.Write(s)
					if err != nil {
						log.Println(err)
					}
				}
				continue
			}
			value, ok := users.Load(m.To)
			if ok {
				conn := value.(*websocket.Conn)
				log.Printf("send: %+v\n", m)
				_, err = conn.Write(msg[:n])
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func main() {
	serverPtr := flag.String("server", ":1234", "--server=:1234")
	pathPtr := flag.String("path", "/echo", "--path=/echo")
	log.Printf("server: %s, path: %s\n", *serverPtr, *pathPtr)
	http.Handle(*pathPtr, websocket.Handler(handle))
	err := http.ListenAndServe(*serverPtr, nil)
	if err != nil {
		log.Panicln(err)
	}
}
