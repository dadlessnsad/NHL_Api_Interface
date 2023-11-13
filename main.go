package main

import (
	"log"
	"net/http"
	"nhl_interface/routes"
	"nhl_interface/services"
	"os"

	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	_ = godotenv.Load()
	r := mux.NewRouter()

	lmt := tollbooth.NewLimiter(50, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	routes.HandleRoutes(r)
	lmtHandler := tollbooth.LimitFuncHandler(lmt, func(w http.ResponseWriter, req *http.Request) {
		r.ServeHTTP(w, req)
	})

	chainedHandler := c.Handler(lmtHandler)
	chainedHandler = handlers.LoggingHandler(os.Stdout, chainedHandler)
	chainedHandler = handlers.RecoveryHandler()(chainedHandler)

	services.ConnectDB()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
		services.CloseDB()
		os.Exit(1)
	}
	log.Println("Server running on port " + port)

	err := http.ListenAndServe(":"+port, chainedHandler)
	if err != nil {
		services.CloseDB()
		os.Exit(1)
	}

}
