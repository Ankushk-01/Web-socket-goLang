package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	// clientlist map[*Client]bool
	clients ClientList
	sync.RWMutex
	Handlers map[string]EventHandler // stores the name and eventHandler function of it
}

//	func newManager() manager {
//		return manager{}
//	}
func newManager() *Manager { // (*) value of the variable and (&) address of the variable
	m := &Manager{
		clients:  make(ClientList), //so that not get null pointer exception
		Handlers: make(map[string]EventHandler),
	}
	m.setupEventHandler()
	return m
}

func (m *Manager) setupEventHandler() {
	m.Handlers[eventSendMessage] = sendMessage
}

func sendMessage(event Event, c *Client) error {
	fmt.Println("event : ", event)
	return nil
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.Handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("There is not such type of Event")
	}
}

var webSocketUpgrader = websocket.Upgrader{
	CheckOrigin:    handleOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (m *Manager) handleWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new Connection")
	// upgrade the regular http request to Http 101 into websocket
	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Websocket is not applied with error %v", err)
		return
	}

	client := newClient(conn, m)

	m.addClient(client)
	// Start Client messages

	go client.readMessages()
	go client.writeMessages()
	// conn.Close()
}

func (m *Manager) addClient(c *Client) {
	m.Lock()
	m.clients[c] = true
	defer m.Unlock()
}
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
	defer m.Unlock()
}

func handleOrigin(r *http.Request) bool {

	origin := r.Header.Get("Origin")
	switch(origin){
	case "localhost://8080":{
		return true
	}
	}
	return false
}
