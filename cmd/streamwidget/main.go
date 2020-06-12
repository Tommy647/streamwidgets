package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello, Cath, how are you?")
	http.HandleFunc("/", HelloServer)
	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}

// HelloServer is just a test function
func HelloServer(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
	}
	log.Println(r.URL.Path[1:])
}
