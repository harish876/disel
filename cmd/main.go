package main

import (
	"encoding/json"
	"flag"
	"strings"
	"time"

	disel "github.com/harish876/disel/core"
)

type ExampleBody struct {
	Foo string `json:"foo"`
}

func main() {
	flag.Parse()

	host := "0.0.0.0"
	port := 42069

	app := disel.New()
	// app.UseThreadPool(1)
	app.Log.SetLevel(disel.INFO).Build()

	app.GET("/", func(c *disel.Context) error {
		return c.Status(200).Send("Success")
	})

	app.GET("/echo", func(c *disel.Context) error {
		if len(c.Request.PathParams) > 0 {
			content := strings.Join(c.Request.PathParams, "/")
			return c.Status(200).Send(content)
		} else {
			return c.Status(200).Send("Success")
		}
	})

	app.POST("/echo", func(c *disel.Context) error {
		var body ExampleBody
		time.Sleep(5 * time.Second)
		if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
			return c.Status(400).Send("Unable to Decode Body")
		}
		app.Log.Info("Request Foo from Body ", body.Foo)
		return c.Status(200).JSON(body)
	})

	app.Log.Infof("Starting Server... on Port %d\n", port)
	app.Log.Fatal(app.ServeHTTP(host, port))
}
