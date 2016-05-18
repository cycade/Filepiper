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
    "math/rand"
    "strings"
)

// 本机数据库地址
const dbURL = "127.0.0.1:27017"
// 提取码本 理论上四位提取码可对应62*62*62*62 = 14776336个文件
const seed = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"


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
// 建立关于上传文件信息的结构
type metaFile struct{
	Ecode string
	Filename string
	MD5 string
	UploadDate time.Time
	DownloadTimes int
	IsValid bool
}


// 通用的检查错误函数，在一般情况下使用
func check(err error) {
	if err != nil {
		panic(err)
	}
}
// 显示首页欢迎画面，可以去掉了
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "Hello! This is a amazing Skydrive!")
}


// 随机生成四位 string 格式提取码
func ecodeGenerator() string {
	randomNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
	extractCode := ""
	for i := 0; i < 4; i++ {
		extractCode += string(seed[randomNumber.Intn(61)])
	}
	return extractCode
}
// 生成上传文件的metaInfo
func metaInfoGenerator(Filename, MD5 string, UploadDate time.Time) metaFile {
	metaInfo := metaFile{ecodeGenerator(), Filename, MD5, UploadDate,
		0, true}
	return metaInfo
}
// 将生成的metaInfo填入到metafiles数据表中
func metaInfotoDB(metaInfo metaFile) {
	session, err := mgo.Dial(dbURL)
	check(err)
	defer session.Close()
	metafiles := session.DB("Filepiper").C("metafiles")
	result := metaFile{}
	err = metafiles.Find(bson.M{"Ecode": metaInfo.Ecode}).One(&result)
	if err != nil {
		err = metafiles.Insert(&metaInfo)
	} else {
		metaInfo.Ecode = ecodeGenerator()
		metaInfotoDB(metaInfo)
	}
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
		gfs := session.DB("Filepiper").GridFS("fs")
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
		fileUploadDate := fileConfig.UploadDate()
		defer fileConfig.Close()

		// 调用函数生成提取码对应表
		metaInfo := metaInfoGenerator(UploadFileName, fileMD5, fileUploadDate)
		metaInfotoDB(metaInfo)

		// 在客户端输出提取码
		fmt.Fprintf(w, "The ExtractCode is %s", metaInfo.Ecode)
	}
}


// 检查提取码中的每个字符是否在seed中
func ContainsAll(seed, ecode string) bool {
	judgement := ""
	for i := 0; i < 4 ; i++ {
		if strings.Contains(seed, string(ecode[i])) {
			judgement += "t"
		} else {
			judgement += "f"
		}
	}

	if strings.Contains(judgement, "f") {
		return false
	} else {
		return true
	}
}
// 在ContainsAll的基础上检查提取码的长度是否合法
func checkEcode(ecode string) bool {
	if len(ecode) == 4 && ContainsAll(seed, ecode) {
		return true
	} else {
		return false
	}
}

func metaInfoFind(ecode string) (metaFile, error) {
	session, err := mgo.Dial(dbURL)
	check(err)
	defer session.Close()

	metafiles := session.DB("Filepiper").C("metafiles")
	result := metaFile{}
	err = metafiles.Find(bson.M{"ecode": ecode}).One(&result)
	return result, err
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var uploadTemplate = template.Must(template.ParseFiles("update.html"))
		err := uploadTemplate.Execute(w, nil)
		check(err)
	} else {
		r.ParseForm()
		ecode := r.Form["extractCode"][0]
		if checkEcode(ecode) == false {
			fmt.Fprintf(w, "You have to put a exactly ExtractCode!")
			return
		}
		// 连接mongodb服务器
		session, err := mgo.Dial(dbURL)
		check(err)
		defer session.Close()

		metaInfo, err := metaInfoFind(ecode)
		if err != nil {
			fmt.Fprintf(w, "The extractCode is Invalid or the file does not exists.")
			return
		}

		gfs := session.DB("Filepiper").GridFS("fs")
		result := gfsFile{}
		err = gfs.Find(bson.M{"md5": metaInfo.MD5}).One(&result)
		if err != nil {
			fmt.Fprintf(w, "The file does not exist.")
			return
		} else {
			//fmt.Fprintf(w, "The file is downloading. Please wait!")
			fmt.Println(result)

			var downloadFile *mgo.GridFile
		    downloadFile, err = gfs.OpenId(result.Id)
		    check(err)
		    defer downloadFile.Close()
		    http.ServeContent(w, r, metaInfo.Filename, downloadFile.UploadDate(), downloadFile)
		}
	}
}



// 测试提取文件
// UUJL and md5 is cf192aa9baccc274293bcb1b162a5ffb
// mmQk and md5 is cf192aa9baccc274293bcb1b162a5ffb
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
