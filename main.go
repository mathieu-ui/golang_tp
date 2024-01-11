package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("111")

	f := make(Ferme)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/moutonlist", func(w http.ResponseWriter, r *http.Request) {
		for sheep := range f {
			a := "Nom: " + f[sheep].Name + ", age: " + strconv.Itoa(f[sheep].Age) + ", poid: " + strconv.FormatFloat(f[sheep].Weight, 'f', 2, 64) + "\n"
			fmt.Println(a)
			w.Write([]byte(a))
		}
	})

	r.Post("/mouton", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var m Sheep

		err := json.NewDecoder(r.Body).Decode(&m)

		fmt.Println()
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(m)

		AddSheep(f, m.Id, m.Id, m.Name, m.Age, m.Weight)

		// b, err := io.ReadAll(r.Body)
		// if err == nil {
		// 	fmt.Println(b)
		// 	myString := string(b[:])
		// 	fmt.Println(myString)
		// 	fmt.Println(myString[12])
		// }

	})

	http.ListenAndServe(":3333", r)

}

type Ferme map[int]Sheep

type Sheep struct {
	Id     int
	Name   string
	Age    int
	Weight float64
}

func NewSheep(idSheep int, nameSheep string, ageSheep int, weightSheep float64) Sheep {
	nSheep := Sheep{
		Id:     idSheep,
		Name:   nameSheep,
		Age:    ageSheep,
		Weight: weightSheep}
	return nSheep
}

func AddSheep(f1 Ferme, cle int, id int, nom string, age int, poid float64) {
	f1[cle] = NewSheep(id, nom, age, poid)
}
