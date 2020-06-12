package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello, Cath, how are you?")
	http.HandleFunc("/", FolderHandler)
	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}

// FolderHandler handles folder paths
func FolderHandler(w http.ResponseWriter, r *http.Request) {
	folder := r.URL.Path[1:]
	file := fmt.Sprintf("./examples/%s/index.html", folder) // @todo: just testing code

	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = fmt.Fprint(w, string(data))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
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
