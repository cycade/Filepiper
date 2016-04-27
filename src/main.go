package main

import (
	"fmt"
	"html/template"
	"log"
	"io/ioutil"
	"net/http"
	"strings"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key: ", k)
		fmt.Println("val: ", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello cycade!")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		fmt.Println("username: ", r.Form["username"])
		fmt.Println("password: ", r.Form["password"])
	}
}

var uploadTemplate = template.Must(template.ParseFiles("update.html"))

func indexHandle(w http.ResponseWriter, r *http.Request) {
	if err := uploadTemplate.Execute(w, nil); err != nil {
		log.Fatal("Execute: ", err.Error())
		return
	}
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Fatal("FormFile: ", err.Error())
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("Close: ", err.Error())
			return
		}
	}()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("ReadAll: ", err.Error())
		return
	}

	w.Write(bytes)
}

func main() {
	http.HandleFunc("/", Welcome)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", indexHandle)
	http.HandleFunc("/upload/upload", uploadHandle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}