package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketController struct{}

func (ws WebsocketController) Start(w http.ResponseWriter, r *http.Request) {
	var dialer = websocket.Dialer{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   "ipaddress:port",
			Path:   "/",
		}),
	}
	dialer.Dial("", r.Header)
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}
