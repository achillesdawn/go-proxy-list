package proxylist

import (
	"fmt"
	"io"
	"log"
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
