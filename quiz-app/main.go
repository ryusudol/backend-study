package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()

	newStore, err := NewNoteStore("notes.db")
	if err != nil {
		log.Fatal(err)
	}
	srv := &server{store: newStore}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /notes", srv.listNotes)
	mux.HandleFunc("POST /notes", srv.createNote)
	mux.HandleFunc("GET /notes/{id}", srv.getNote)
	mux.HandleFunc("PUT /notes/{id}", srv.updateNote)
	mux.HandleFunc("DELETE /notes/{id}", srv.deleteNote)

	mux.HandleFunc("/notes", methodNotAllowed("GET, POST"))
	mux.HandleFunc("/notes/{id}", methodNotAllowed("GET, PUT, DELETE"))
	mux.HandleFunc("/", notFound)

	log.Println("listening on", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
