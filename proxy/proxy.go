package proxytest

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func writeToFile(res *http.Response) {
	file, err := os.Create("result.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	n, err := io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("wrote", n, "bytes")
}

func makeRequest(targetUrl string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	reader := io.LimitReader(bufio.NewReader(res.Body), 984735+100_000)
	b, err := io.ReadAll(reader)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return b, nil
}
