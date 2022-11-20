package message

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Message struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func RunServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			newMessage := Message{}
			b, err := ioutil.ReadAll(r.Body)

			err2 := json.Unmarshal(b, &newMessage)
			if err2 != nil {
				log.Println(err2)
			}
			if err != nil {
				panic(err)
			}

		})

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Println("starting server at", addr)
	server.ListenAndServe()
}
