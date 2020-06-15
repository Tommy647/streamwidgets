package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Starting up")
	http.HandleFunc("/", FolderHandler("examples"))
	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}

// FolderHandler handles folder paths
func FolderHandler(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		path := r.URL.Path[1:]
		file := fmt.Sprintf("%s/%s/%s/index.html", dir, folder, path) // @todo: just testing code
		log.Printf("using: %s", file)
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
