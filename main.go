package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"path"
	"io"
	"bytes"
	"strings"
)

// downloaded images storage file
const downloadDir string = "img/"

var imageUrlList []string
var num chan bool = make(chan bool, 5)

func main() {
	url := "https://github.com/daimajia?tab=followers"
	getImage(url)
}

func getImage(url string) (err error) {

	InitFolder()

	imageUrlList = getImageUrl(url)

	// download the image.
	for _, imgUrl := range imageUrlList {
		fmt.Println(imgUrl)
		num <- true
		go getImg(imgUrl)
	}

	// return nil when all tasks were finished
	return nil
}

// download the images with url
func getImg(url string) (n int64, err error) {
	// use the userId to name the img
	reg := regexp.MustCompile("/\\d{7,9}")
	urlId := reg.FindString(url)
	paths := strings.Split(urlId, "/")
	var name string
	if len(paths) > 1 {
		name = paths[len(paths)-1]
	}
	fmt.Println(name)
	out, err := os.Create("img/" + string(name) + ".jpg")
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	n, err = io.Copy(out, bytes.NewReader(pix))
	<-num
	return
}

// get image url
func getImageUrl(url string) (imageUrls []string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}

	// regexp
	imgUrlMatcher := regexp.MustCompile("[a-zA-z]+://[^\\s]*[s=]\\d{3}")

	imgUrls := imgUrlMatcher.FindAllString(string(body), -1)
	return imgUrls
}

func InitFolder() (err error) {
	// create the storage folder
	if err := os.Mkdir(path.Dir(downloadDir), os.ModePerm); err != nil {
		fmt.Println("file is already exisist!")
		return err
	}
	fmt.Println("Storage folder was created.")
	return nil
}
