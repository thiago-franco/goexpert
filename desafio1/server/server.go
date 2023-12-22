package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CurrencyExchange struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

const API_TIMEOUT = 200 * time.Millisecond

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), API_TIMEOUT)
	defer cancel()
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro no servidor ao fazer requisicao: %v\n", err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro no servidor ao receber resposta: %v\n", err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Response error: %v\n", err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	var data CurrencyExchange
	err = json.Unmarshal(body, &data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unmarshal error: %v\n", err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(data)

	name := data.Usdbrl.Name
	bid := data.Usdbrl.Bid

	db, err := criarDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	defer db.Close()

	inserirCotacao(db, name, bid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(bid))
}

func criarDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "desafio1.db")
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Erro ao criar banco")
	}

	createTableSQL := `
        CREATE TABLE IF NOT EXISTS cotacoes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            bid TEXT
        );
    `
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Erro ao criar tabela")
	}
	return db, nil
}

func inserirCotacao(db *sql.DB, name string, bid string) {
	insertDataSQL := `
        INSERT INTO cotacoes (name, bid) VALUES (?, ?);
    `
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := db.ExecContext(ctx, insertDataSQL, name, bid)
	if err != nil {
		fmt.Println(err)
		return
	}
}
