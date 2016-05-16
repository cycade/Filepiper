package main

import (
	"fmt"
	"html/template"
	"log"
	"io"
	"net/http"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
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
// 通用的检查错误函数，在一般情况下使用
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "Hello! This is a amazing Skydrive!")
}

// 随机生成四位 string 格式提取码
func ecodeProducer() string {
	var seed = [62]string{"A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R",
		"S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c",
		"d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y",
		"z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

	randomNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
	extractCode := ""
	for i := 0; i < 4; i++ {
		extractCode += seed[randomNumber.Intn(61)]
	}
	return extractCode
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 如果用GET方式请求，则渲染update.html文件
		var uploadTemplate = template.Must(template.ParseFiles("update.html"))
		err := uploadTemplate.Execute(w, nil)
		check(err)
	} else {
		// 如果用POST方式请求，则解析文件并存入数据库中
		file, handler, err := r.FormFile("file") // 解析文件
		if err != nil {
			fmt.Fprintf(w, "You should choose a file to UPLOAD.")
			return
		}
		UploadFileName := handler.Filename // 获取上传文件名
		defer file.Close()
		// 连接数据库服务器
		session, err := mgo.Dial(dbURL)
		check(err)
		defer session.Close()
		// 创建一个类型为 *GridFile 的文件
		gfs := session.DB("FPusers").GridFS("Files")
		gfile, err := gfs.Create(UploadFileName)
		check(err)
		// io.Copy 将file的内容复制到gfile中
		_, err = io.Copy(gfile, file)
		check(err)
		gfile.Close()

		// 将新建的 *GridFile 文件打开并获md5值
		fileConfig, err := gfs.Open(UploadFileName)
		check(err)
		fileMD5 := fileConfig.MD5()
		defer fileConfig.Close()
		
		// 新建
		ecode := ecodeProducer()
		
		// 在客户端输出提取码
		fmt.Fprintf(w, "The ExtractCode is %s", fileMD5) 
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var uploadTemplate = template.Must(template.ParseFiles("update.html"))
		err := uploadTemplate.Execute(w, nil)
		check(err)
	} else {
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
}

// 测试提取文件
// cf192aa9baccc274293bcb1b162a5ffb

func main() {
	fmt.Println("Server starting.")
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/404", http.NotFound)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
