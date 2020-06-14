package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func download(URL string, ch chan<- []byte) {
	resp, err := http.Get(URL)
	check(err)

	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	ch <- body
}

func main() {
	now := time.Now()
	ch := make(chan []byte)

	file, err := ioutil.ReadFile("ceps.txt")
	check(err)

	ceps := strings.Split(string(file), "\r\n")

	for _, cep := range ceps {
		URL := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
		go download(URL, ch)
	}

	for _, cep := range ceps {
		err := ioutil.WriteFile("./db/"+cep+".json", <-ch, 0644)
		check(err)
	}

	fmt.Printf("Durou %2.f segundos", time.Since(now).Seconds())
}
