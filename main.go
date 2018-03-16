package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/capitalone/fpe/ff1"
	"github.com/gorilla/mux"
)

var (
	port        = flag.String("port", "8080", "fpeService port")
	keyString   = flag.String("key", "FF4359D8D580AA4F7F036D6F04FC6A94", "key for the FF1 algorithm")
	tweakString = flag.String("tweak", "D8E7920AFA330A73", "tweak for the FF1 algorithm")
)

type req struct {
	words    string
	respChan chan string
}

type Server struct {
	c   chan req
	FF1 ff1.Cipher
}

func NewServer(key, tweak []byte) *Server {
	FF1, err := ff1.NewCipher(62, 8, key, tweak)
	if err != nil {
		panic(err)
	}
	s := Server{
		c:   make(chan req),
		FF1: FF1,
	}
	return &s
}

func (s *Server) encrypter() {
	var err error
	for {
		request := <-s.c
		log.Println(request)
		originalSplit := strings.Split(request.words, " ")
		tokenised := make([]string, len(originalSplit))
		for i, word := range originalSplit {
			tokenised[i], err = s.FF1.Encrypt(word)
			if err != nil {
				log.Println("failed to encrypt", word)
				panic(err) // TODO handle error gracefully
			}
		}
		request.respChan <- strings.Join(tokenised, " ")
	}
}

func (s *Server) encryptHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	ret := make(chan string)
	log.Println(string(body))
	s.c <- req{string(body), ret}
	fmt.Fprintf(w, <-ret)
}

func (s *Server) NewRouter() *mux.Router {
	type Route struct {
		Name        string
		Pattern     string
		Method      string
		HandlerFunc http.HandlerFunc
	}
	routes := []Route{
		Route{
			"Encyrpt",
			"/encrypt",
			"GET",
			s.encryptHandler,
		},
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

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

	s := NewServer(key, tweak)
	r := s.NewRouter()

	go s.encrypter()

	http.Handle("/", r)

	log.Println("serving on", *port)
	err = http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Panicf(err.Error())
	}

}
