package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /cotacao", cotacao)
	log.Println("Server initiated at port 8080")
	http.ListenAndServe(":8080", nil)
}

func cotacao(w http.ResponseWriter, r *http.Request) {
	log.Println("/cotacao request received")
	defer log.Println("/cotacao request completed")

	w.Write([]byte("Hello World"))
}
