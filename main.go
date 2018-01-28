package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
  "github.com/kataras/iris/middleware/recover"
  "github.com/parnurzeal/gorequest"
  "encoding/json"
  "fmt"
)

type Response struct {
  Ip string `json:"ip"`
  Region string `json:"region"`
  City string `json:"city"`
  Country string `json:"country"`
  Org string `json:"org"`
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
  app.Use(logger.New())
  
  // set the view engine target to ./templates folder
	app.RegisterView(iris.HTML("./templates", ".html").Reload(true))

	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
    ip := ctx.FormValue("ip")
    res := Response{}
    url := ""

    if ip == "" {
      url = "https://ipinfo.io/json"
    } else {
      url = "https://ipinfo.io/" + ip + "/json"
    }

    resp, body, errs := gorequest.New().Get(url).End()

    if errs != nil {
      ctx.JSON(iris.Map{"message": "error"}) 
    }

    err := json.Unmarshal([]byte(body), &res)

    if err != nil {
      ctx.JSON(iris.Map{"message": "error"}) 
    }

    s := &res

    fmt.Print(resp)

    ctx.ViewData("title", s.Ip)
		ctx.ViewData("location", s.Country + " " + s.Region + " " + s.City)
		ctx.ViewData("org", s.Org)
		// same file, just to keep things simple.
		if err := ctx.View("index.html"); err != nil {
			ctx.Application().Logger().Infof(err.Error())
		}
	})

	// http://localhost:8080
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}