package disel

import (
	"context"
	"encoding/json"
	"io"
)

type Context struct {
	Request  HttpRequest
	Ctx      context.Context
	Response HttpResponse
}

func (c *Context) ReadBody() (string, error) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *Context) Status(status int) *Context {
	c.Response.status = status
	return c
}

func (c *Context) ContentType(contentType string) *Context {
	c.Response.contentType = contentType
	return c
}

func (c *Context) Send(body string) error {
	if c.Response.status == 0 {
		c.Response.status = 200
	}
	if c.Response.contentType == "" {
		c.Response.contentType = "text/plain"
	}
	c.Response.body = Serialize(
		c.Response.status,
		body,
		c.Response.contentType,
	)
	return nil
}

func (c *Context) JSON(body any) error {
	mBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	c.Response.contentType = "application/json"
	return c.Send(string(mBody))
}
