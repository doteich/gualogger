package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"gualogger/logging"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	Port     int    `mapstruct:"port"`
	Endpoint string `mapstruct:"endpoint"`
	Username string `mapstruct:"username"`
	Password string `mapstruct:"password"`
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
	pongDeadline = 10 * time.Second
	pingInterval = (pongDeadline * 9) / 10
)

type manager struct {
	sync.RWMutex
	clients map[*client]bool
	secret  string
}

type client struct {
	addedTS    time.Time
	connection *websocket.Conn
	manager    *manager
}

type inbound_event struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

func (ws *Websocket) Initialize(ctx context.Context) error {

	str := fmt.Sprintf("%s:%s", ws.Username, ws.Password)

	sec := base64.StdEncoding.EncodeToString([]byte(str))

	ws.manager = manager{clients: make(map[*client]bool), secret: sec}

	_, err := url.Parse(ws.Endpoint)

	if err != nil {
		return err
	}

	http.HandleFunc(ws.Endpoint, ws.upgrade)

	go ws.manager.verifyClients()

	go startServer(ws.Port)

	return nil
}

func (ws *Websocket) Publish(ctx context.Context, p Payload) error {

	var e error

	for c, auth := range ws.manager.clients {

		if !auth {
			continue
		}
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

func (ws *Websocket) upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		logging.Logger.Error(fmt.Sprintf("error while upgrading websocket connection: %s", err.Error()), "func", "websocket_upgrade")
		return
	}

	ws.manager.Lock()

	defer ws.manager.Unlock()

	c := client{connection: conn, manager: &ws.manager, addedTS: time.Now()}

	ws.manager.clients[&c] = false

	go c.readMessages()
	go c.writeMessages()
}

func startServer(port int) {
	addr := fmt.Sprintf(":%d", port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logging.Logger.Error(fmt.Sprintf("unable to start websocket server on port %d: %s", port, err.Error()), "func", "websocket_startServer")
		os.Exit(100)
	}
}

func (c *client) readMessages() {

	defer c.manager.removeClient(c)

	if err := c.connection.SetReadDeadline(time.Now().Add(pongDeadline)); err != nil {
		return
	}

	c.connection.SetPongHandler(c.pongHandler)

	for {
		var inb inbound_event
		if err := c.connection.ReadJSON(&inb); err != nil {
			logging.Logger.Warn(fmt.Sprintf("received malformed websocket message - removing client: %s", err.Error()), "func", "websocket_readmessages")
			return
		}

		switch inb.Name {
		case "authentication_message":
			if inb.Payload == c.manager.secret {
				c.manager.authenticateClient(c)
				continue
			}
			return

		default:
			return
		}

	}
}

func (c *client) writeMessages() {

	tick := time.NewTicker(pingInterval)
	defer func() {
		tick.Stop()
		c.manager.removeClient(c)
	}()

	for {
		select {

		case <-tick.C:
			// Send the Ping
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				logging.Logger.Warn(fmt.Sprintf("unable to write ping message: %s", err.Error()), "func", "websocket_writemessages")
				return
			}
		}

	}
}

func (m *manager) authenticateClient(c *client) {
	_, ok := m.clients[c]

	if ok {
		m.Lock()
		m.clients[c] = true
		m.Unlock()
	}

}

func (m *manager) removeClient(c *client) {
	m.Lock()
	defer m.Unlock()

	_, ok := m.clients[c]

	if ok {
		c.connection.WriteMessage(websocket.CloseMessage, []byte(""))
		c.connection.Close()
		delete(m.clients, c)

		return
	}
	//logging.Logger.Warn("unable to remove websocket client from map", "func", "websocket_removeclient")
}

func (m *manager) verifyClients() {
	for {
		for c, auth := range m.clients {

			if !auth && time.Since(c.addedTS) > 10*time.Second {
				logging.Logger.Info("websocket client was unable to authenticate in the specified timeframe - removing client", "func", "websocket_verifyclients")
				m.removeClient(c)
			}

		}

		time.Sleep(30 * time.Second)
	}
}

func (c *client) pongHandler(pongMsg string) error {

	return c.connection.SetReadDeadline(time.Now().Add(pongDeadline))
}
