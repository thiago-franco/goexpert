package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	TIMEOUT = 300 * time.Millisecond
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, res.Body)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	if res.StatusCode == http.StatusOK {
		f.Write([]byte(fmt.Sprintf("Dólar: %s", body)))
		return
	}
}
