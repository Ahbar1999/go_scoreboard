package main

import (
	"fmt"
	"net/http"
	"log"
	"strings"
	"strconv"
	"encoding/json"
	"math/rand"
	"errors"
)

type Player struct {
	Id int			`json:"id"`
	Name string		`json:"name"`
	Country string	`json:"country"`
	Score int		`json:"score"`
}

var players = make([]Player, 0)

func newPlayer(name string, country string, score int) (*Player, error) {
	if len(name) == 0 || len(name) > 15 {
		return nil, errors.New("length of name should be between 1 and 15")
	} else if len(country) != 2 {
		return nil, errors.New("country code invalid")
	}
	
	players = append(players, Player{len(players), name, country, score})
	return &players[len(players) - 1], nil 
}


func getPlayers(w http.ResponseWriter, r *http.Request) {
	urlString := r.URL.String()

	switch true {
		case strings.HasPrefix(urlString, "/players/rank/"):
			i := strings.Index(urlString, ":")
			rank, err := strconv.Atoi(urlString[i + 1:])
			if err != nil || rank < 0 || rank > len(players) {
				// bad request
				w.WriteHeader(http.StatusBadRequest)	
				return
			}

			w.WriteHeader(http.StatusOK)
			resp, err := json.Marshal(players[rank])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))	
				return
			}
			w.Write(resp)

		case strings.HasPrefix(urlString, "/players/random"):
			w.WriteHeader(http.StatusOK)
			resp, err := json.Marshal(players[rand.Intn(len(players))])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))	
				return
			}
			w.Write(resp)

		case urlString == "/players":
			w.WriteHeader(http.StatusOK)
			resp, err := json.Marshal(players)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))	
				return
			}
			w.Write(resp)
		
		default:
			w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	for i := 0; i < 10; i += 1{
		newPlayer("aa", "IN", i * 100)
		fmt.Println(players)	
	}

	mux := http.NewServeMux()

	// register handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.String()))
	})	

	mux.Handle("/players/", http.HandlerFunc(getPlayers))
	
	log.Fatal(http.ListenAndServe(":" + "8080", mux))
}

