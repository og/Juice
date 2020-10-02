package juice

import (
	ogjson "github.com/og/json"
	gconv "github.com/og/x/conv"
	ge "github.com/og/x/error"
	gtest "github.com/og/x/test"
	"log"
	"testing"
)
type ReqHome struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
}
type ReplyHome struct {
	IDNameAge string `json:"idNameAge"`
}
func StartServe(addr string) {
	serve := NewServe(ServeOption{
		OnCatchError: func(c *Context, errInterface interface{}) error {
			log.Print(errInterface)
			return nil
		},
	})
	serve.Action(POST, "/", func(c *Context) (reject error) {
		req := ReqHome{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		reply := ReplyHome{}
		reply.IDNameAge  = req.ID + ":" + req.Name + ":" + gconv.IntString(req.Age)
		return c.Bytes(ogjson.Bytes(reply))
	})
	ge.Check(serve.Listen(addr))
}
func TestTest(t *testing.T) {
	as := gtest.NewAS(t)
	addr := ":1111"
	go StartServe(addr)
	resp := HttpTest(addr, HttpTestRequestJSON(POST, "/", ReqHome{
		ID: "a",
		Name: "b",
		Age: 3,
	}))
	// resp.String resp.Bytes resp.BindJSON resp.ExpectJSON 任选其一即可
	{
		as.Equal(resp.String(), `{"idNameAge":"a:b:3"}`)
	}
	{
		reply := ReplyHome{}
		resp.BindJSON(&reply)
		as.Equal(reply, ReplyHome{
			IDNameAge: "a:b:3",
		})
	}
	{
		resp.ExpectJSON(t, ReplyHome{IDNameAge: "a:b:3",})
	}
}
