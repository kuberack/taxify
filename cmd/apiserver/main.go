package main

import (
	"fmt"
	"log"
	"net/http"

	"kuberack.com/taxify/internal/api"
)

func main() {

	h, err := api.NewServerWithMiddleware()

	if err != nil {
		log.Fatal(err.Error())
	}

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())

	fmt.Printf("help")
}
