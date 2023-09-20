## Ping 

Server send a special method to the client(Web Browser) in our case to see that the connection is on or not we can set the time interval of it 

## Pong

the pong is the response of the client to the server that the connection is on or not and it is implicit by the browser so no need to implement it on server


## Jumble frames

jumble frames to limit incoming message size 

`c.connection.SetReadLimit(512)`

other way is to set the origin that is set the specifiy the trafic on the port 

```
func handleOrigin(r *http.Request) bool {

	origin := r.Header.Get("Origin")
	switch(origin){
	case "localhost://8080":{
		return true
	}
	}
	return false
}
```