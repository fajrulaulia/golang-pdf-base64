package gopdfbase64

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func AddSignHandler(r *mux.Router) {
	r.HandleFunc("/", IndexController).Methods("GET")
	r.HandleFunc("/process", ProcessController).Methods("POST")

}

func IndexController(w http.ResponseWriter, r *http.Request) {
	var data = map[string]string{
		"Title":    "Convert PDF <> base64",
		"Subtitle": "Easy to convert PDF <> base64",
	}
	var t, err = template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, data)
}

func ProcessController(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(10); err != nil {
		log.Println("Error on r.ParseMultipartForm", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	if handler.Header.Get("Content-Type") != "application/pdf" {

		var res = map[string]string{
			"Request": "no-valid-pdf",
		}
		t, err := template.ParseFiles("output.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t.Execute(w, res)
		return

	}

	data, err := ioutil.ReadAll(uploadedFile)
	if err != nil {
		log.Println(err)
		return
	}
	result := ByteToBase64(data, handler.Filename)

	var res = map[string]string{
		"Output":  result,
		"Request": "OK",
	}

	t, err := template.ParseFiles("output.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, res)

}

var lock sync.Mutex

func Load(data []byte, filename string) {
	lock.Lock()
	defer lock.Unlock()

	dirr, err := os.Getwd()
	if err != nil {
		log.Println("err", err)
	}

	err = CreateFile(data, dirr+"/files/"+filename)
	log.Println("Save your data in local with filename", filename)
	if err != nil {
		log.Println("Error when createfile by spies", err)
		return
	}

	log.Println("Add filename into repos.txt")

	f, err := os.OpenFile("repos.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(filename + "\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func ByteToBase64(data []byte, filename string) string {
	go Load(data, filename+time.Now().String()+".pdf")
	log.Println("Converting file " + filename + "...")
	str := base64.StdEncoding.EncodeToString(data)
	log.Println("Convert success to save", filename)
	return str
}

func WriteBase64ToFile(data, filename string) error {
	dec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}

	return nil

}

func CreateFile(fil []byte, filename string) error {

	err := ioutil.WriteFile(filename, fil, 0644)
	if err != nil {
		return err
	}

	return nil
}
