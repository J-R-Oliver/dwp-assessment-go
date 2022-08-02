package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	u := fmt.Sprintf("http://127.0.0.1:%s/health", os.Getenv("PORT"))

	r, err := http.Get(u) //nolint:gosec
	if err != nil {
		os.Exit(1)
	}

	r.Body.Close()

	if r.StatusCode != http.StatusNoContent {
		os.Exit(1)
	}
}
