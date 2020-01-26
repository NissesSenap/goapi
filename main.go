package main

import (
	"fmt"
	"net/http"
	"os"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Asset not found \n"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Running API v1\n"))
}

func main() {
	http.HandleFunc("/", rootHandler)
	err := http.ListenAndServe("localhost:11111", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
