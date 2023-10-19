package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

type Book struct {
	Title string `json:"title"`
	Autor string `json:"author"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

func Server() {

	var books []Book

	mux := http.NewServeMux()

	tmpl := template.Must(template.ParseFiles("./public/index.html"))

	data, err := os.ReadFile("books.json")

	if err != nil {
		panic("error ao converter")
	}

	err = json.Unmarshal(data, &books)

	if err != nil {
		panic("erro ao dar unmarshal")
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, books)
	})

	http.ListenAndServe(":8080", mux)

}
