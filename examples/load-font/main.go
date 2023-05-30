package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/pchchv/captcha"
	"golang.org/x/image/font/gofont/goregular"
)

func indexHandle(w http.ResponseWriter, _ *http.Request) {
	doc, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	doc.Execute(w, nil)
}

func captchaHandle(w http.ResponseWriter, _ *http.Request) {
	img, err := captcha.New(150, 50, func(options *captcha.Options) {
		options.FontScale = 0.8
	})
	if err != nil {
		fmt.Fprint(w, nil)
		fmt.Println(err.Error())
		return
	}

	img.WriteImage(w)
}

func main() {
	if err := captcha.LoadFont(goregular.TTF); err != nil {
		panic(err)
	}

	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/captcha", captchaHandle)

	fmt.Println("Server start at port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
