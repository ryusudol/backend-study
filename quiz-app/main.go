package main

import (
	"log"
	"net/http"
)

func main() {
	srv := &server{store: NewNoteStore()}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /notes", srv.listNotes)
	mux.HandleFunc("POST /notes", srv.createNote)
	mux.HandleFunc("GET /notes/{id}", srv.getNote)
	mux.HandleFunc("PUT /notes/{id}", srv.updateNote)
	mux.HandleFunc("DELETE /notes/{id}", srv.deleteNote)

	mux.HandleFunc("/notes", methodNotAllowed("GET, POST"))
	mux.HandleFunc("/notes/{id}", methodNotAllowed("GET, PUT, DELETE"))
	mux.HandleFunc("/", notFound)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
