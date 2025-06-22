package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	server := NewServer()

	http.HandleFunc("/", server.loginHandler)
	http.HandleFunc("/register", server.registerHandler)
	http.HandleFunc("/verify", server.verifyHandler)
	http.HandleFunc("/setup", server.setupHandler)
	http.HandleFunc("/success", server.successHandler)
	http.HandleFunc("/qr", server.qrHandler)

	fmt.Println("サーバーを開始します: http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
