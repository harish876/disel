package disel

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/common-nighthawk/go-figure"
)

type Disel struct {
	Options      map[string]string
	Log          *Logger
	GetHandlers  RadixTree
	PostHandlers RadixTree
}

type DiselHandlerFunc func(c *Context) error

func New() Disel {
	return Disel{
		Options:      make(map[string]string),
		Log:          InitLogger(),
		GetHandlers:  NewRadixTree(),
		PostHandlers: NewRadixTree(),
	}
}

func (d *Disel) AddOption(optionKey string, optionValue string) {
	d.Options[optionKey] = optionValue
}

func (d *Disel) GET(path string, handler DiselHandlerFunc) error {
	d.Log.Debug("Registered GET Route for", path)
	d.GetHandlers.Insert(path, &handler)
	return nil
}

func (d *Disel) POST(path string, handler DiselHandlerFunc) error {
	d.Log.Debug("Registered POST Route for", path)
	d.PostHandlers.Insert(path, &handler)
	return nil
}

func (d *Disel) ServeHTTP(host string, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("Failed to bind to port")
		os.Exit(1)
	}

	welcomeMsg := figure.NewColorFigure("Disel", "", "green", true)
	welcomeMsg.Print()
	d.Log.Infof("Server Stated on port... %d", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go d.handleConnection(conn)
	}
}

func (d *Disel) execHandler(ctx *Context) error {
	var handler DiselHandlerFunc
	if ctx.Request.Method == "GET" {
		node, found := d.GetHandlers.Search(ctx.Request.Path)
		d.Log.Debug("Incoming GET Route Path is", ctx.Request.Path)
		if !found {
			handler = nil
		} else {
			handler = *node.Value.(*DiselHandlerFunc)
			d.Log.Debug("GET Handler is", handler)
		}

	} else if ctx.Request.Method == "POST" {
		node, found := d.PostHandlers.Search(ctx.Request.Path)
		d.Log.Debug("Incoming POST Route Path is", ctx.Request.Path)
		if !found {
			handler = nil
		} else {
			handler = *node.Value.(*DiselHandlerFunc)
			d.Log.Debug("POST Handler is", handler)
		}
	} else {
		handler = nil
	}
	if handler == nil {
		ctx.Status(404).Send(fmt.Sprintf("Route Not found for Incoming Path %s", ctx.Request.Path))
		return nil
	}
	_, cancel := context.WithTimeout(ctx.Ctx, time.Second*10)
	defer cancel()
	if err := handler(ctx); err != nil {
		ctx.Status(http.StatusInternalServerError).Send("Not Found")
		return err
	}
	log.Println(ctx.Response.body)
	return nil
}

func (d *Disel) handleConnection(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		recievedBytes, err := conn.Read(buf)
		if err == io.EOF || err != nil {
			log.Println(err)
			break
		}
		request := buf[:recievedBytes]
		rawRequest := string(request)
		parsedRequest := DeserializeRequest(rawRequest)
		d.Log.Debug("Raw Request is", rawRequest)
		ctx := &Context{
			Request: parsedRequest,
			Ctx:     context.Background(),
		}

		_ = d.execHandler(ctx)
		sentBytes, err := conn.Write([]byte(ctx.Response.body))
		if err != nil {
			d.Log.Debug("Error writing response: ", err.Error())
		}
		d.Log.Debug("Sent Bytes to Client: ", sentBytes)
	}
}
