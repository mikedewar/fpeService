package main

import (
	"encoding/hex"
	"flag"
	"log"
	"net/http"
)

var (
	port        = flag.String("port", "8080", "fpeService port")
	keyString   = flag.String("key", "FF4359D8D580AA4F7F036D6F04FC6A94", "key for the FF1 algorithm")
	tweakString = flag.String("tweak", "D8E7920AFA330A73", "tweak for the FF1 algorithm")
)

func main() {

	flag.Parse()
	key, err := hex.DecodeString(*keyString)
	if err != nil {
		panic(err)
	}
	tweak, err := hex.DecodeString(*tweakString)
	if err != nil {
		panic(err)
	}

	s := NewFPEServer(key, tweak)
	r := s.NewRouter()

	go s.encrypter()

	http.Handle("/", r)

	log.Println("serving on", *port)
	err = http.ListenAndServeTLS(":"+*port, "cert.pem", "key.pem", nil)
	if err != nil {
		log.Panicf(err.Error())
	}

}
