package main

import (
	"io"
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

	ctx := r.Context()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
