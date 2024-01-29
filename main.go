package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "db"
)

func main() {
	fmt.Println("Server is running on port 3333")

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Erreur lors de la connexion à la base de données:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Erreur lors de la vérification de la connexion à la base de données:", err)
		return
	}

	fmt.Println("Connexion à la base de données réussie")

	db_Query := "CREATE TABLE IF NOT EXISTS moutons (id SERIAL PRIMARY KEY, name TEXT NOT NULL, age FLOAT NOT NULL, weight FLOAT NOT NULL);"

	_, err = db.Exec(db_Query)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête:", err)
		return
	}

	fmt.Println("Requête exécutée avec succès")

	insertQuery := "INSERT INTO moutons (name, age, weight) VALUES ('juju', 10, 100)"

	_, err = db.Exec(insertQuery)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête d'ajout:", err)
		return
	}

	fmt.Println("Requête d'ajout exécutée avec succès")


	insertQuery1 := "SELECT * FROM moutons"
	rows, err := db.Query(insertQuery1)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête de sélection:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age float64
		var weight float64

		err := rows.Scan(&id, &name, &age, &weight)
		if err != nil {
			fmt.Println("Erreur lors de la lecture des données:", err)
			return
		}

		fmt.Printf("ID: %d, Name: %s, Age: %f, Weight: %f\n", id, name, age, weight)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Erreur lors de la récupération des lignes:", err)
		return
	}


	f := make(Ferme)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bergerie is Open"))
	})

	r.Get("/moutonlist", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT name, age, weight FROM moutons")
		if err != nil {
			fmt.Println("Erreur lors de l'exécution de la requête de sélection:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			var age float64
			var weight float64

			err := rows.Scan(&name, &age, &weight)
			if err != nil {
				fmt.Println("Erreur lors de la lecture des données:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			a := "Nom: " + name + ", age: " + strconv.FormatFloat(age, 'f', 2, 64) + ", poids: " + strconv.FormatFloat(weight, 'f', 2, 64) + "\n"
			fmt.Println(a)
			w.Write([]byte(a))
		}
	})

	r.Post("/mouton", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var m Sheep

		err := json.NewDecoder(r.Body).Decode(&m)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(m)

		AddSheep(f, m.Id, m.Id, m.Name, m.Age, m.Weight)

	})

	r.Post("/dellmouton", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var m Sheep

		err := json.NewDecoder(r.Body).Decode(&m)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		delete(f, m.Id)
	})

	r.Post("/updatemouton", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var m Sheep

		err := json.NewDecoder(r.Body).Decode(&m)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, ok := f[m.Id]; ok {
			AddSheep(f, m.Id, m.Id, m.Name, m.Age, m.Weight)
		} else {
			fmt.Println("Mouton non trouvé")
			http.Error(w, "Mouton non trouvé", http.StatusNotFound)
		}
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
