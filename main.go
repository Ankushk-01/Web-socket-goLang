package main

import (
	"log"
	"net/http"
)

func main() {
	setApi()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Connot Serve at port : 8080 and eror is : %v", err)
	}
	log.Println("Serving at port : 8080")

}

func setApi() {
	manager  := newManager()
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws",manager.handleWS)
}
