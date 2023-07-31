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
	"sort"
)

type Player struct {
	Id int			`json:"id"`
	Name string		`json:"name"`
	Country string	`json:"country"`
	Score int		`json:"score"`
}

var players = make([]*Player, 0)
var newId = 0

func newPlayer(name string, country string, score int) (*Player, error) {
	if len(name) == 0 || len(name) > 15 {
		return nil, errors.New("length of name should be between 1 and 15")
	} else if len(country) != 2 {
		return nil, errors.New("country code invalid")
	}
	
	players = append(players, &Player{newId, name, country, score})
	newId += 1
	return players[len(players) - 1], nil 
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	urlString := r.URL.String()
	fmt.Println(urlString)
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})	
	switch r.Method {
	case "GET":
		switch true {
		case strings.HasPrefix(urlString, "/players/rank/"):
			i := strings.Index(urlString, ":")
			
			rank, err := strconv.Atoi(urlString[i + 1:])
			if err != nil || rank < 0 || rank >= len(players) {
				// bad request
				w.WriteHeader(http.StatusBadRequest)
				if err != nil {
					w.Write([]byte(err.Error()))
				} else {
					w.Write([]byte("invalid rank"))
				}	
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

		case strings.HasPrefix(urlString, "/players"):
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
	
	case "DELETE":
		i := strings.Index(urlString, ":")
		deleted := false	
		id, _ := strconv.Atoi(urlString[i + 1:])
		for i, p := range players {
			if p.Id == id {
				// shrink slice
				if i == len(players) - 1 {
					players = players[:len(players) - 1]
				} else {	
					players = append(players[:i], players[i + 1:]...)
				}
				deleted = true
			}
		}
		
		if !deleted {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("Player with id: : %v NOT FOUND", id))) 
			return	
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Successfully deleted player with id: %v", id))) 
		}
	
	case "POST":
		r.ParseForm()
		score, err := strconv.Atoi(r.Form.Get("score"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}

		newPlayer(r.Form.Get("name"), r.Form.Get("country"), score)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully added a new player"))
	
	case "PUT":
		r.ParseForm()
		i := strings.Index(urlString, ":")
		if i == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Send a request with ':id'"))
			return
		}
		modified := false	
		id, _ := strconv.Atoi(urlString[i + 1:])
		for i, p := range players {
			if p.Id == id {
				
				if r.Form.Has("name") {
					players[i].Name = r.Form.Get("name")
				} 
				if r.Form.Has("score") {
					score, _ := strconv.Atoi(r.Form.Get("score")) 
					players[i].Score = score
				} 
				if r.Form.Has("name") || r.Form.Has("score") {
					modified = true
				}
				break
			}
		}
		
		if !modified {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("Player with id: : %v NOT FOUND", id))) 
			return	
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Successfully modified player with id: %v", id))) 
		}
	}	
}

func main() {
	for i := 0; i < 10; i += 1{
		newPlayer("aa", "IN", i * 100)
		fmt.Println(players)	
	}

	mux := http.NewServeMux()

	// register handlers
	mux.Handle("api/players/", http.HandlerFunc(getPlayers))
	
	log.Fatal(http.ListenAndServe(":" + "8080", mux))
}

