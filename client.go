package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	exchangeRate, err := getServerAnswer()
	if err != nil {
		log.Println(err)
		return
	}

	extended_exchange := "Dólar: " + exchangeRate + "\n"
	err = saveExchangeRatesInfoInFile(extended_exchange)
	if err != nil {
		log.Println(err)
	}

	io.Copy(os.Stdout, strings.NewReader(extended_exchange))
}

func getServerAnswer() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8080/cotacao", nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("deu ruim: ", res.StatusCode)
		return "", errors.New("deu ruim: " + res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(resBody), nil
}

func saveExchangeRatesInfoInFile(exchangeRate string) error {
	filename := "cotacao.txt"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(exchangeRate); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
