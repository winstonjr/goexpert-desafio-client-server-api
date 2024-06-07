package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
)

func main() {
	db, err := maybeCreateSQLLiteDatabase()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /cotacao", Cotacao{dbConn: db})

	log.Println("Server initiated at port 8080")
	http.ListenAndServe(":8080", mux)
}

type Cotacao struct {
	dbConn *sql.DB
}

func (cotacao Cotacao) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("/cotacao request received")
	defer log.Println("/cotacao request completed")
	if cotacao.dbConn != nil {
		log.Println("Conexão com banco de dados estabelecida")
	}

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

func maybeCreateSQLLiteDatabase() (*sql.DB, error) {
	const create string = `CREATE TABLE IF NOT EXISTS exchange_rates (id INTEGER NOT NULL PRIMARY KEY, codein TEXT NOT NULL,
							  name TEXT NOT NULL, high TEXT NOT NULL, low TEXT NOT NULL, var_bid TEXT NOT NULL,
							  pct_change  TEXT NOT NULL, bid TEXT NOT NULL, ask TEXT NOT NULL, timestamp TEXT NOT NULL,
							  create_date TEXT NOT NULL);`

	db, err := sql.Open("sqlite3", "exchange_rates.db")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	log.Println("connected to database exchange_rates.db")
	return db, nil
}

//func fileExists(filePath string) (bool, error) {
//	info, err := os.Stat(filePath)
//	if err == nil {
//		return !info.IsDir(), nil
//	}
//	if errors.Is(err, os.ErrNotExist) {
//		return false, nil
//	}
//	return false, err
//}
