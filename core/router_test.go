package disel

import (
	"fmt"
	"testing"
)

func TestRadixTree(t *testing.T) {
	handler := func(c *Context) error {
		return c.Status(200).Send("Success")
	}
	fmt.Println("\n---- Radix tree test ------ ")
	rt := NewRadixTree()
	rt.Insert("/", handler)
	rt.Insert("/echo", handler)
	rt.Insert("/user-agent", handler)
	rt.Insert("/files", handler)
	rt.Insert("/test", handler)
	rt.PrintAllWords()

}
