package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type DolarPriceResult struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, error := io.ReadAll(res.Body)
	if error != nil {
		panic(error)
	}

	var dp DolarPriceResult

	json.Unmarshal(body, &dp)
	writeDolarPrice(dp)

	io.Copy(os.Stdout, res.Body)
}

func writeDolarPrice(dp DolarPriceResult) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString("Dólar: " + dp.Bid)
	if err != nil {
		panic(err)
	}
	defer f.Close()

}
