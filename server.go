package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DolarPrice struct {
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

type DolarPriceResult struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
	}

	dp, err := getDolarPrice()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = saveDolarPrice(dp); err != nil {
		log.Println(err)
	}

	dpr := DolarPriceResult{Bid: dp.Usdbrl.Bid}

	json.NewEncoder(w).Encode(dpr)
}

func getDolarPrice() (*DolarPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var dp DolarPrice
	if err = json.NewDecoder(res.Body).Decode(&dp); err != nil {
		return nil, err
	}

	return &dp, nil
}

func saveDolarPrice(dp *DolarPrice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	log.Println("Dolar price: ", dp.Usdbrl.Bid)
	db, err := getDBConnection()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO dollarPrice (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, createDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, dp.Usdbrl.Code, dp.Usdbrl.Codein, dp.Usdbrl.Name, dp.Usdbrl.High, dp.Usdbrl.Low, dp.Usdbrl.VarBid, dp.Usdbrl.PctChange, dp.Usdbrl.Bid, dp.Usdbrl.Ask, dp.Usdbrl.Timestamp, dp.Usdbrl.CreateDate)
	if err != nil {
		return err
	}
	return nil
}

func getDBConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		return nil, err
	}
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS dollarPrice (code string, codein string, name string, high string, low string, varBid string, pctChange string, bid string, ask string, timestamp string, createDate string)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	return db, nil
}
