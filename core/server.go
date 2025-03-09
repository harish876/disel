package disel

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/common-nighthawk/go-figure"
)

type ThreadPoolOptions struct {
	pool          ThreadPool
	useThreadPool bool
}
type DiselOptions struct {
	threadPool ThreadPoolOptions
}

func NewDiselOptions() DiselOptions {
	return DiselOptions{
		threadPool: ThreadPoolOptions{
			pool:          ThreadPool{},
			useThreadPool: false,
		},
	}
}

type Disel struct {
	Options        DiselOptions
	Log            *Logger
	GetHandlers    RadixTree
	PostHandlers   RadixTree
	PutHandlers    RadixTree
	DeleteHandlers RadixTree
}

type DiselHandlerFunc func(c *Context) error

func New() Disel {
	return Disel{
		Options:        DiselOptions{},
		Log:            InitLogger(),
		GetHandlers:    NewRadixTree(),
		PostHandlers:   NewRadixTree(),
		PutHandlers:    NewRadixTree(),
		DeleteHandlers: NewRadixTree(),
	}
}

func (d *Disel) UseThreadPool(poolSize int) {
	d.Options.threadPool.useThreadPool = true
	var wg sync.WaitGroup
	d.Options.threadPool.pool = NewThreadPool(poolSize, &wg)

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
	displayWelcomeMessage()
	d.Log.Infof("Server Stated on port... %d", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}
		defer conn.Close()
		if d.Options.threadPool.useThreadPool {
			d.Options.threadPool.pool.Add(func() {
				d.handleConnection(conn)
			})
		} else {
			go d.handleConnection(conn)
		}
	}
	d.Options.threadPool.pool.Wait()
	os.Exit(1)
	return nil
}

func displayWelcomeMessage() {
	welcomeMsg := figure.NewColorFigure("Disel", "", "green", false)
	welcomeMsg.Print()
	fmt.Printf("\n")
}

func (d *Disel) execHandler(ctx *Context) error {
	var handler DiselHandlerFunc
	switch ctx.Request.Method {

	case "GET":
		node, found := d.GetHandlers.Search(ctx.Request.Path)
		d.Log.Debug("Incoming GET Route Path is", ctx.Request.Path)
		if !found {
			handler = nil
		} else {
			handler = *node.Value.(*DiselHandlerFunc)
			d.Log.Debug("GET Handler is", handler)
		}

	case "POST":
		node, found := d.PostHandlers.Search(ctx.Request.Path)
		d.Log.Debug("Incoming POST Route Path is", ctx.Request.Path)
		if !found {
			handler = nil
		} else {
			handler = *node.Value.(*DiselHandlerFunc)
			d.Log.Debug("POST Handler is", handler)
		}

	case "PUT", "PATCH":
		node, found := d.PutHandlers.Search(ctx.Request.Path)
		d.Log.Debug("Incoming PUT Route Path is", ctx.Request.Path)
		if !found {
			handler = nil
		} else {
			handler = *node.Value.(*DiselHandlerFunc)
			d.Log.Debug("PUT Handler is", handler)
		}

	case "DELETE":
		node, found := d.DeleteHandlers.Search(ctx.Request.Path)
		d.Log.Debug("Incoming DELETE Route Path is", ctx.Request.Path)
		if !found {
			handler = nil
		} else {
			handler = *node.Value.(*DiselHandlerFunc)
			d.Log.Debug("DELETE Handler is", handler)
		}
	default:
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
	return nil
}

func (d *Disel) handleConnection(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		recievedBytes, err := conn.Read(buf)
		if err == io.EOF || err != nil {
			d.Log.Debug(err)
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

		err = d.execHandler(ctx)
		if err != nil {
			d.Log.Error(err)
		}
		sentBytes, err := conn.Write([]byte(ctx.Response.body))
		if err != nil {
			d.Log.Debug("Error writing response: ", err.Error())
		}
		d.Log.Debug("Sent Bytes to Client: ", sentBytes)
	}
}
