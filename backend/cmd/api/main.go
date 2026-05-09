package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akolesnov/football58/backend/internal/config"
)

func main() {
	cfg := config.Load()

	http.HandleFunc("/", hello)

	if err := http.ListenAndServe(cfg.HTTPAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "hello world")
}
