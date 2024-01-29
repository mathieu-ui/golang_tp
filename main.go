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

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bergerie is Open"))
	})

	r.Get("/moutonlist", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, age, weight FROM moutons")
		if err != nil {
			fmt.Println("Erreur lors de l'exécution de la requête de sélection:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			a := "ID: " + strconv.Itoa(id) + ", Nom: " + name + ", age: " + strconv.FormatFloat(age, 'f', 2, 64) + ", poids: " + strconv.FormatFloat(weight, 'f', 2, 64) + "\n"
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

		insertQuery := "INSERT INTO moutons (name, age, weight) VALUES ($1, $2, $3)"
		_, err = db.Exec(insertQuery, m.Name, m.Age, m.Weight)
		if err != nil {
			fmt.Println("Erreur lors de l'exécution de la requête d'ajout:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Requête d'ajout exécutée avec succès")

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

		deleteQuery := "DELETE FROM moutons WHERE name = $1 AND age = $2 AND weight = $3 AND id = $4"
		
		_, err = db.Exec(deleteQuery, m.Name, m.Age, m.Weight, m.Id)
		if err != nil {
		    fmt.Println("Erreur lors de l'exécution de la requête de suppression:", err)
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		    return
		}

		fmt.Println("Mouton supprimé avec succès")
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

		updateQuery := "UPDATE moutons SET name = $1, age = $2, weight = $3 WHERE id = $4"
		_, err = db.Exec(updateQuery, m.Name, m.Age, m.Weight, m.Id)
		if err != nil {
			fmt.Println("Erreur lors de l'exécution de la requête de mise à jour:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Mouton mis à jour avec succès")
	})

	http.ListenAndServe(":3333", r)

}

type Sheep struct {
	Id     int
	Name   string
	Age    int
	Weight float64
}
