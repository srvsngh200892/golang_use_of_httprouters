
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to furit world!\n")
}

type FruitList struct {

	Id   string `json:"id"`
	Name  string `json:"name"`

}

type JsonResponse struct {
	// Reserved field to add some meta information to the API response
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}


// This acts as the fruits storage
var fruitstore = make(map[string]*FruitList)

// Handler for get all food
// GET /fruits
func FruitIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fruits := []*FruitList{}
	for _, fruit := range fruitstore {
		fruits = append(fruits, fruit)
	}
	response := &JsonResponse{Data: &fruits}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

// Handler for the get one paticaular food
// GET /fruits/:id
func Show(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	fruit, ok := fruitstore[id]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if !ok {
		// fruit not found
		w.WriteHeader(http.StatusNotFound)
		response := JsonErrorResponse{Error: &ApiError{Status: 404, Title: "Fruit Not Found"}}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}
	response := JsonResponse{Data: fruit}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func CreateFruit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var fruit FruitList
	decoder.Decode(&fruit)
	fruitstore["4"] = &FruitList{
		Id:   fruit.Id,
		Name:  fruit.Name,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fruits := []*FruitList{}
	for _, fruit := range fruitstore {
		fruits = append(fruits, fruit)
	}
	response := &JsonResponse{Data: &fruits}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/fruits", FruitIndex)
	router.GET("/fruits/:id", Show)
	router.POST("/fruits", CreateFruit)

	fruitstore["1"] = &FruitList{
		Id:   "1",
		Name:  "Apple",
	}

	fruitstore["2"] = &FruitList{
		Id:   "2",
		Name:  "Mango",
	}

	fruitstore["3"] = &FruitList{
		Id:   "3",
		Name:  "orange",
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
