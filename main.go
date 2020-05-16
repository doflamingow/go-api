package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type heroes struct {
	ID     int    `from:"ID" json:"id"`
	Name   string `from:"Name" json:"name"`
	Title  string `from:"Title" json:"title"`
	Armor  int    `from:"Armor" json:"armor"`
	Damage int    `from:"Damage" json:"damage"`
	Hp     int    `from:"Hp" json:"hp"`
}

type responseWrapper struct {
	Status string   `from:"Status" json:"status"`
	Data   []heroes `from"Data" json:"data"`
}

type handleSuccess struct {
	Status string `from:"Status" json:"status"`
	Data   heroes `from"Data" json:"data"`
}

type handlerError struct {
	Status string `from:"Status" json:"status"`
	Error  string `from:"Error" json:"error"`
}

var data []heroes

func getAllHeroes(w http.ResponseWriter, r *http.Request) {
	var datas = responseWrapper{
		Status: "success",
		Data:   data,
	}

	// var dat []heroes

	// dat = append(dat, heroes{ID: 1, Name: "Wiro Sableng", Title: "The Warrior", Armor: 999, Damage: 999, Hp: 999,})

	//heroes{ID: 1, Name: "Wiro Sableng", Title: "The Warrior", Armor: 999, Damage: 999, Hp: 999}

	json.NewEncoder(w).Encode(datas)
}

func getHeroesById(w http.ResponseWriter, r *http.Request) {

	getID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		respons := handlerError{
			Status: "fail",
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(respons)
		return
	}

	if getID > len(data) || getID <= 0 {
		respons := handlerError{
			Status: "fail",
			Error:  "data dosnt exists",
		}
		json.NewEncoder(w).Encode(respons)
		return
	}

	for i := range data {
		if getID == data[i].ID {
			var datas = handleSuccess{
				Status: "success",
				Data:   data[i],
			}
			json.NewEncoder(w).Encode(datas)
			return
		}
	}

	respons := handlerError{
		Status: "fail",
		Error:  "id dosnt exists",
	}
	json.NewEncoder(w).Encode(respons)
}

func addHeroes(w http.ResponseWriter, r *http.Request) {
	var structHeroes heroes

	json.NewDecoder(r.Body).Decode(&structHeroes)

	if structHeroes.ID == 0 || structHeroes.Name == "" || structHeroes.Title == "" || structHeroes.Armor == 0 || structHeroes.Damage == 0 || structHeroes.ID == 0 {

		var datas = handlerError{
			Status: "fail",
			Error:  "please fill all body",
		}
		json.NewEncoder(w).Encode(datas)
		return
	}

	data = append(data, structHeroes)

	// result := handleSuccess{
	//     Status: "success",
	//     Data:   structHeroes,
	// }
	json.NewEncoder(w).Encode(structHeroes)
}

func deleteHeroes(w http.ResponseWriter, r *http.Request) {
	var isExists bool
	var tmpData []heroes

	getID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		respons := handlerError{
			Status: "fail",
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(respons)
		return
	}

	if getID > len(data) || getID <= 0 {
		respons := handlerError{
			Status: "fail",
			Error:  "data dosnt exists",
		}
		json.NewEncoder(w).Encode(respons)
		return
	}

	for i := range data {
		if getID != data[i].ID {
			tmpData = append(tmpData, data[i])
		} else if getID == data[i].ID {
			isExists = true
		}
	}

	if !isExists {
		respons := handlerError{
			Status: "fail",
			Error:  "id dosnt exists",
		}
		json.NewEncoder(w).Encode(respons)
		return
	}

	data = tmpData

	respons := responseWrapper{
		Status: "success",
		Data:   data,
	}
	json.NewEncoder(w).Encode(respons)
}

func main() {
	port := "8080"

	data = append(data, heroes{ID: 1, Name: "Balmond", Title: "Tank", Armor: 199, Damage: 999, Hp: 999})
	data = append(data, heroes{ID: 2, Name: "Hayabusa", Title: "Assasins", Armor: 399, Damage: 999, Hp: 999})
	data = append(data, heroes{ID: 3, Name: "Fanny", Title: "Assasins", Armor: 120, Damage: 43, Hp: 100})
	data = append(data, heroes{ID: 4, Name: "Layla", Title: "Marksman", Armor: 270, Damage: 50, Hp: 100})
	data = append(data, heroes{ID: 4, Name: "Jhonson", Title: "Tank", Armor: 280, Damage: 50, Hp: 100})
	data = append(data, heroes{ID: 4, Name: "Hylos", Title: "Tank", Armor: 250, Damage: 50, Hp: 100})
	data = append(data, heroes{ID: 4, Name: "Granger", Title: "Marksman", Armor: 220, Damage: 50, Hp: 100})


	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/heroes", logger(getAllHeroes)).Methods(http.MethodGet)
	router.HandleFunc("/heroes/{id}", logger(getHeroesById)).Methods(http.MethodGet)
	router.HandleFunc("/heroes", logger(addHeroes)).Methods(http.MethodPost)
	router.HandleFunc("/heroes/{id}", logger(deleteHeroes)).Methods(http.MethodDelete)

	log.Println("Server starting at http://localhost" + port)
	http.ListenAndServe(":"+port, router)
}

func logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("URL [%v]: %v dipanggil pada jam %v \n", r.Method, r.URL.Path, time.Now().Format("13-Mei-2020 15:04:05"))
		next.ServeHTTP(w, r)
	})
}
