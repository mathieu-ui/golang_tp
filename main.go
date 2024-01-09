package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("111")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/mouton", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		// var m Sheep

		//err := json.NewDecoder(r.Body).Decode(&m)

		//fmt.Println()
		// if err != nil {
		// 	fmt.Println(err)
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }
		// fmt.Fprintf(w, "Person: %+v", m)

		b, err := io.ReadAll(r.Body)
		if err == nil {
			fmt.Println(b)
			myString := string(b[:])
			fmt.Println(myString)
			fmt.Println(myString[12])
		}

	})

	http.ListenAndServe(":3333", r)

}

type Sheep struct {
	Id     int
	Name   string
	Age    int
	Weight float64
}
