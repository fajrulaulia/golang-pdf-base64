package main

import (
	"log"
	"net/http"
	"os"

	c "go-pdf-base64/app"

	"github.com/gorilla/mux"
)

func main() {
	// dat, err := os.ReadFile("./test.pdf")
	// if err != nil {
	// 	log.Println("error on os.ReadFile=", err)
	// 	return
	// }
	// log.Println("c.ByteToBase64(dat)", c.ByteToBase64(dat))

	r := mux.NewRouter()

	c.AddSignHandler(r)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Running on port", port)

	http.ListenAndServe(":"+port, r)
}
