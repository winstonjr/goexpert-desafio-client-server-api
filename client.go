package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8080/cotacao", nil)
	if err != nil {
		log.Println(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}
