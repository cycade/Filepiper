package main

import (
	"fmt"
	"html/template"
	"log"
	"io"
	"net/http"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    "crypto/md5"
)

const dbURL = "127.0.0.1:27017"

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
		session, err := mgo.Dial(dbURL)
		check(err)
		result := Users{}
		where := session.DB("FPusers").C("Users")
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
		session, err := mgo.Dial(dbURL)
		check(err)
		resultin := Users{}
		where := session.DB("FPusers").C("Users")
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
	err := uploadTemplate.Execute(w, nil)
	check(err)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "You should choose a file to UPLOAD.")
	}
	UploadFileName := handler.Filename // 获取上传文件名
	defer file.Close()

	session, err := mgo.Dial(dbURL)
	check(err)
	defer session.Close()
	
	db := session.DB("FPusers").GridFS("Files")
	dataset, err := db.Open(UploadFileName)
	if err != nil {
		dataset, err = db.Create(UploadFileName)
		check(err)
	}

	_, err = io.Copy(dataset, file)
	check(err)
	defer dataset.Close()

	fmt.Fprintf(w, "The file has been uploaded and md5 is and the md5 is %s", dataset.Id())
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

