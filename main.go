package main

import (
	"log"
	"net/http"
	"nhl_interface/routes"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	// apply middleware
	chainedHandler := handlers.CORS()(c.Handler(r))
	routes.HandleRoutes(r)

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	log.Println("Listening on port " + port)

	err = http.ListenAndServe(":"+port, chainedHandler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		os.Exit(1)
	}
}
