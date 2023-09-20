package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool
type Client struct {
	// Every client have a one- one connection with web socket
	connection *websocket.Conn
	manager *Manager

	// egress is used to avoid concurrent writes on the web socket connection and it is of type channel
	egress chan Event
}
var (
	pongWaitTime = time.Second *10

	pingInterval = (pongWaitTime*9)/10
)

func newClient(conn *websocket.Conn, mann *Manager) (*Client){
	return &Client{
		connection: conn,
		manager: mann,
		egress: make(chan Event),
	}

}

func (c *Client) readMessages(){
	defer func(){ // when the loops ends then the defer function will executed
		// cleanup the client connection
		c.manager.removeClient(c)
	}()

	c.connection.SetReadLimit(512)  // jumble frames to limit incoming message size 
	// pong message 
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWaitTime));err!=nil{
		log.Println("error in setting pong time ",err)
	}
	c.connection.SetPongHandler(c.pongHandler)
	for{
		_,payload,err := c.connection.ReadMessage();
		if err!= nil{
			if(websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway,websocket.CloseAbnormalClosure)){
				log.Printf("Error reading message : %v",err)
			}
			break
		}
		// writing on egress quick hack
		
		var request Event

		if err := json.Unmarshal(payload,&request);err!=nil{
			log.Println("Error in unmarshalling the payload ",err)
		}
		if err := c.manager.routeEvent(request,c);err!=nil{
			log.Println("error in routeEvent method ",err)
		}

		// log.Println(connectionType) 
		// log.Println(string(payload))  // because payload is in binary form so convert it to string
		// fmt.Println("payload : ",string(payload))
	}
}

func (c *Client) writeMessages(){
	ticker := time.NewTicker(pingInterval)
	defer func ()  {
		c.connection.Close()
	}()

	for{
		select{
		case message,ok:= <-c.egress:  // ok means the egress is fine 
		if !ok{
			if err := c.connection.WriteMessage(websocket.CloseMessage,nil);err!= nil{
				log.Println("Cannot send message to client : ",err)
			}
			return
		}
		data , err:= json.Marshal(message)
		if err!=nil {
			log.Println("Error in marshalling the message ",err)
		}
		if err := c.connection.WriteMessage(websocket.TextMessage,data);err!=nil{
			log.Println("cannot write messages : ",err)
		}
		log.Println("Message Send")

	case <-ticker.C:
		log.Println("Ping")
		if err := c.connection.WriteMessage(websocket.PingMessage,[]byte(``));err!=nil{
			log.Println("error in ping msg ",err)
		}
		}
	}
}

func (c *Client) pongHandler(pong string) error{
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWaitTime))
	
}