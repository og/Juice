package juice

import (
	ogjson "github.com/og/json"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
	serve *Serve
}
func NewContext (w http.ResponseWriter, r *http.Request, serve *Serve) *Context {
	return &Context{
		W: w,
		R: r,
		serve: serve,
	}
}
func (c *Context) Bytes(b []byte) error {
	_, err := c.W.Write(b)
	if err != nil {
		_, err := c.W.Write([]byte("c.Bytes(b) error"))
		if err != nil {panic(err)}
	}
	return nil
}
func (c *Context) JSON(v interface{}) error {
	jsonb, err := ogjson.BytesWithErr(v)
	if err != nil {
		return err
	}
	return c.Bytes(jsonb)
}
func (c *Context) BindRequest(ptr interface{}) error {
	return BindRequest(ptr, c.R)
}

func (c *Context) CheckError(errInterface interface{}) {
	err := c.serve.OnCatchError(c, errInterface)
	if err != nil {
		panic(err)
	}
}