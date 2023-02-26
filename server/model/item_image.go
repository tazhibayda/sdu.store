package model

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Image struct {
	ID     int64  `json:"id"`
	ItemID int64  `json:"itemID"`
	Name   string `json:"name"`
	Data   []byte `json:"data"`
}

func UploadImage(w http.ResponseWriter, r *http.Request) {

	fmt.Println("File Upload")
	err := r.ParseMultipartForm(10 << 30)
	if err != nil {
		return
	}

	file, handler, err := r.FormFile("image")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {

		}
	}(tempFile)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return
	}

}

func ShowImages(w http.ResponseWriter, r *http.Request) {
	tm, _ := template.ParseFiles("templates/test/image.gohtml")
	files, err := ioutil.ReadDir("images")
	if err != nil {
		log.Fatal(err)
	}
	imgs := make([]string, 0)
	for _, file := range files {
		imgs = append(imgs, file.Name())
	}
	err = tm.Execute(w, imgs)
	if err != nil {
		return
	}
}
