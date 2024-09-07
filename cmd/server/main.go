package main

import (
	"groupietracker/pkg/handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	// Configuration des routes et des handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/home", handlers.HomePage)
	mux.HandleFunc("/artist/", handlers.ArtistDetailsPage)
	mux.HandleFunc("/", handlers.NotFoundHandler)
	mux.HandleFunc("/search", handlers.SearchHandler)
	mux.HandleFunc("/suggestions", handlers.SuggestionHandler)

	// Configuration du serveur
	server := &http.Server{
		Addr:              ":8080",           // Adresse du serveur
		Handler:           mux,               // Liste des handlers
		ReadHeaderTimeout: 10 * time.Second,  // Temps autorisé pour lire les headers
		WriteTimeout:      10 * time.Second,  // Temps maximum d'écriture de la réponse
		IdleTimeout:       120 * time.Second, // Temps maximum entre deux requêtes
		MaxHeaderBytes:    1 << 20,           // 1 MB, maximum de bytes pour les headers
	}

	// Lancement du serveur
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
