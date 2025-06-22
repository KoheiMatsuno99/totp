package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	server := NewServer()

	http.HandleFunc("/", server.loginHandler)
	http.HandleFunc("/verify", server.verifyHandler)
	http.HandleFunc("/setup", server.setupHandler)
	http.HandleFunc("/success", server.successHandler)

	fmt.Println("サーバーを開始します: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
