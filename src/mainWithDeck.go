package main

import (
	"fmt"
	"html/template"
	"log"
	//"io"
	"net/http"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    //"os"
)

type Users struct {
	Username string
	Password string
	Email string
	Telephonenumber int
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "Hello! This is a amazing Skydrive!")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		const URL = "127.0.0.1:27017"
		session, err := mgo.Dial(URL)
		if err != nil {
			panic(err)
		}
		result := Users{}
		where := session.DB("FPusers").C("col")
		user := r.Form["username"][0]
		pass := r.Form["password"][0]
		err = where.Find(bson.M{"username": user, "password": pass}).One(&result)
    	if err != nil {
    		fmt.Fprintf(w, "There is no %s or password is wrong.", user)
    	} else {
    		fmt.Fprintf(w, "You are login.")
		}
	}
}


func signupHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("signup.gtpl")
		t.Execute(w, nil)
	} else {
		const URL = "127.0.0.1:27017"
		session, err := mgo.Dial(URL)
		if err != nil {
			panic(err)
		}
		resultin := Users{}
		where := session.DB("FPusers").C("col")
		userInfo := r.Form
		err = where.Find(bson.M{"username": userInfo["username"][0]}).One(&resultin)
		if err != nil {
    		tele, _ := strconv.Atoi(userInfo["telephonenumber"][0])
    		err = where.Insert(&Users{
			userInfo["username"][0], 
			userInfo["password"][0], 
			userInfo["email"][0], 
			tele,
			})
    		if err != nil {
    			panic(err)
    		} else {
    			fmt.Fprintf(w, "signup finished with %s", userInfo["username"][0])
    		}
    	} else {
    		fmt.Fprintf(w, "%s has been existed.", userInfo["username"][0])	
    	}
    }
}



func indexHandler(w http.ResponseWriter, r *http.Request) {
	var uploadTemplate = template.Must(template.ParseFiles("update.html"))
	if err := uploadTemplate.Execute(w, nil); err != nil {
		log.Fatal("Execute: ", err.Error())
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "The file has been uploaded.")
	f, err := db.GridFS("Files").Create("mute.txt")
	check(err)
	n, err := f.Write(file)
	check(err)
	err = f.Close()
	check(err)
}

func main() {
	fmt.Println("Server starting.")
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/upload", indexHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/upload/upload", uploadHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

