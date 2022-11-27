package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ip2location/ip2location-go/v9"
	"laika/ip2geo/utils"
	"net/http"
	"time"
)

type ip2GeoResponse struct {
	CountryLong  string  `json:"country_name"`
	CountryShort string  `json:"country_short_code"`
	City         string  `json:"city"`
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
}

func main() {
	db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB5.BIN")

	if err != nil {
		fmt.Print(err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/ip", func(w http.ResponseWriter, r *http.Request) {
		ip := utils.ReadUserIP(r)
		results, _ := db.Get_all(ip)

		response := ip2GeoResponse{
			CountryLong:  results.Country_long,
			CountryShort: results.Country_short,
			City:         results.City,
			Latitude:     results.Latitude,
			Longitude:    results.Longitude,
		}

		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":2711", r) // because 27 November :)
}
