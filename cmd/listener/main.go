package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const PORT = 6340

func handler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Println(err)
	}

	// simulate work
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	time.Sleep(time.Duration(n) * time.Millisecond)

	w.WriteHeader(http.StatusCreated)

	log.Printf("message: %s\n", b)
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("server started on port: %d", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
