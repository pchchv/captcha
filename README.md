# *captcha* provides an simple, unopinionated API for captcha generation

[![goreference](https://pkg.go.dev/badge/github.com/pchchv/captcha)](https://pkg.go.dev/github.com/pchchv/captcha)
[![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/captcha)](https://goreportcard.com/report/github.com/pchchv/captcha)

## Compatibility

This package uses embed package from Go 1.16.

## Usage

```Go
import "github.com/pchchv/captcha"

func handle(w http.ResponseWriter, r *http.Request) {
	// create a captcha of 150x50px
	data, _ := captcha.New(150, 50)

	// session come from other library such as gorilla/sessions
	session.Values["captcha"] = data.Text
	session.Save(r, w)
	// send image data to client
	data.WriteImage(w)
}

```

## Sample image
![image](examples/captcha.png)


[documentation](https://pkg.go.dev/github.com/pchchv/captcha) |
[example](examples/exa,ple/main.go) |
[font example](examples/load-font/main.go)
