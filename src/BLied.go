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
    "time"
)

const dbURL = "127.0.0.1:27017"
// 建立关于用户信息的数据结构
type Users struct {
	Username string
	Password string
	Email string
	Telephonenumber int
}
// 建立 MongoDB GridFS.Files 对应的结构
type gfsFile struct {
	Id interface{} "_id"
	UploadDate time.Time "uploadDate"
	Length int64 ",minsize"
	MD5 string
	Filename string
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
	file, handler, err := r.FormFile("file") // 解析文件
	if err != nil {
		fmt.Fprintf(w, "You should choose a file to UPLOAD.")
	}
	UploadFileName := handler.Filename // 获取上传文件名
	defer file.Close()

	session, err := mgo.Dial(dbURL) // 连接数据库服务器
	check(err)
	defer session.Close()
	
	gfs := session.DB("FPusers").GridFS("Files") // 返回一个 *GridFS 类型
	gfile, err := gfs.Create(UploadFileName)  // 创建 *GridFile 类型文件
	check(err)

	_, err = io.Copy(gfile, file) // io.Copy 将file的内容复制到gfile中
	check(err)
	gfile.Close()

	fileConfig, err := gfs.Open(UploadFileName) // 将新建的 *GridFile 文件打开
	check(err)
	fileMD5 := fileConfig.MD5() // 获取新文件的md5值
	fileId := fileConfig.Id() // 获取新文件的_id值
	defer fileConfig.Close()

	fmt.Fprintf(w, "The file has been uploaded %s and md5 is %s", fileId, fileMD5) 
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, err := mgo.Dial(dbURL)
	check(err)
	defer session.Close()

	gfs := session.DB("FPusers").GridFS("Files")
	extractCode := r.Form["extractCode"][0]
	result := gfsFile{}
	err = gfs.Find(bson.M{"md5": extractCode}).One(&result)
	if err != nil {
		fmt.Fprintf(w, "The file is not exist.")
	} else {
		//fmt.Fprintf(w, "The file is downloading. Please wait!")
	
	fmt.Println(result)

	var downloadFile *mgo.GridFile
    downloadFile, err = gfs.OpenId(result.Id)
    check(err)
    defer downloadFile.Close()
    http.ServeContent(w, r, result.Filename, downloadFile.UploadDate(), downloadFile)
	}
}



func main() {
	fmt.Println("Server starting.")
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/upload", indexHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/upload/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
