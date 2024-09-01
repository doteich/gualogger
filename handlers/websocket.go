package handlers

import (
	"context"
	"fmt"
	"gualogger/logging"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	Port     int    `mapstruct:"port"`
	Endpoint string `mapstruct:"endpoint"`
	manager  manager
}

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type manager struct {
	sync.RWMutex
	clients map[*client]bool
}

type client struct {
	connection *websocket.Conn
	manager    *manager
}

func (ws *Websocket) Initialize(ctx context.Context) error {

	ws.manager = manager{clients: make(map[*client]bool)}

	_, err := url.Parse(ws.Endpoint)

	if err != nil {
		return err
	}

	http.HandleFunc(ws.Endpoint, ws.upgrade)

	go startServer(ws.Port)

	return nil
}

func (ws *Websocket) upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		logging.Logger.Error(fmt.Sprintf("error while upgrading websocket connection: %s", err.Error()), "func", "websocket_upgrade")
		return
	}

	ws.manager.RWMutex.Lock()

	c := client{connection: conn, manager: &ws.manager}

	ws.manager.clients[&c] = false

}

func startServer(port int) {
	addr := fmt.Sprintf(":%d", port)

	http.ListenAndServe(addr, nil)
}

func (ws *Websocket) Publish(ctx context.Context, p Payload) error {

	var e error

	for c := range ws.manager.clients {

		if err := c.connection.WriteJSON(p); err != nil {
			e = err
			continue
		}
	}

	return e
}

func (ws *Websocket) Shutdown(ctx context.Context) error {
	return nil
}

func (c *client) ReadMessages() {
	for {
		select {}
	}
}
