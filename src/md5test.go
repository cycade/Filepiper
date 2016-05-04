package main

import (
	"fmt"
	"net/http"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

const dbURL = "127.0.0.1:27017"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func seekHandler (w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(dbURL)
	check(err)
	defer session.Close()

	db := session.DB("FPusers") // 返回一个 *Database 类型
	gfs := db.GridFS("Files") // *GridFS 类型

	gfile, err := gfs.Open("23253.png")
	check(err)
	defer gfile.Close()

	fileMD5 := gfile.MD5()
	fileID := gfile.Id()
	// c := session.DB("FPusers").C("Files.files") //返回一个 *Collection 类型
	fmt.Printf("MD5: %s\nID: %s", fileMD5, fileID)
	fmt.Fprintf(w, "I am %s", fileMD5)
}

func main() {
	fmt.Println("Server starting.")
	http.HandleFunc("/", seekHandler)

	err := http.ListenAndServe(":9090", nil)
	check(err)
}