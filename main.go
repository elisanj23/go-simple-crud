package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deleting book...")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range books {

		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)

		}
	}

	// If the book with the given ID was not found, return an error response
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Book with ID %s not found", params["id"])
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
		}
	}

}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		fmt.Fprintf(w, "Error decoding")
	}

	book.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updating book...")

	//set json type
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)

	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(books)

		}
	}
}

func main() {
	r := mux.NewRouter()

	books = append(books,
		Book{
			ID:     "1",
			Title:  "The Alchemist",
			Author: &Author{FirstName: "John", LastName: "Elisan"},
		})

	books = append(books,
		Book{
			ID:     "2",
			Title:  "The Book of Testament",
			Author: &Author{FirstName: "James", LastName: "Elisan"},
		})

	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	port := ":8000"
	fmt.Println("Starting server on port ", port)
	log.Fatal(http.ListenAndServe(port, r))
}
