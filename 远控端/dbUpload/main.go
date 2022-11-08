package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const UPLOAD_PATH string = "C:/"

type Img struct {
	Id     bson.ObjectId `bson:"_id"`
	ImgUrl string        `bson:"imgUrl"`
}

func main() {
	http.HandleFunc("/entrance", Entrance)
	http.HandleFunc("/uploadImg", UploadImg)
	http.ListenAndServe(":8000", nil)
}

func Entrance(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("uploadImg.html")
	t.Execute(w, nil)
}

func UploadImg(w http.ResponseWriter, r *http.Request) {
	var img Img
	img.Id = bson.NewObjectId()

	r.ParseMultipartForm(1024)
	imgFile, imgHead, imgErr := r.FormFile("img")
	if imgErr != nil {
		fmt.Println("1", imgErr)
		return
	}
	defer imgFile.Close()

	imgFormat := strings.Split(imgHead.Filename, ".")
	img.ImgUrl = img.Id.Hex() + "." + imgFormat[len(imgFormat)-1]

	image, err := os.Create(UPLOAD_PATH + img.ImgUrl)
	if err != nil {
		fmt.Println("2", err)
		return
	}
	defer image.Close()

	_, err = io.Copy(image, imgFile)
	if err != nil {
		fmt.Println("3", err)
		return
	}

	// session, err := mgo.Dial("localhost")
	session, err := mgo.Dial("mongodb://127.0.0.1:27017")
	//默认连接本地mongodb
	if err != nil {
		fmt.Println("4", err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	err = session.DB("test").C("test").Insert(img)
	//本地的数据库及表，可能需要自行修改一下，我本地测试的时候用的都是test
	if err != nil {
		fmt.Println("5", err)
		return
	}
	fmt.Println("上传成功")
}
