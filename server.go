package main

import (
	"log"
	"strings"

	"github.com/capitalone/fpe/ff1"
)

type req struct {
	words    string
	respChan chan string
}

type FPEServer struct {
	c   chan req
	FF1 ff1.Cipher
}

func NewFPEServer(key, tweak []byte) *FPEServer {
	FF1, err := ff1.NewCipher(62, 8, key, tweak)
	if err != nil {
		panic(err)
	}
	s := FPEServer{
		c:   make(chan req),
		FF1: FF1,
	}
	return &s
}

func (s *FPEServer) encrypter() {
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
