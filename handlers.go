package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *FPEServer) encryptHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	ret := make(chan string)
	log.Println(string(body))
	s.c <- req{string(body), ret}
	fmt.Fprintf(w, <-ret)
}
