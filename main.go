package main

import (
	"log"
	"net/http"
	"time"

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

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
