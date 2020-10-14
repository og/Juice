package juicetest_test

import (
	ogjson "github.com/og/json"
	"github.com/og/juice"
	juicetest "github.com/og/juice/test"
	gconv "github.com/og/x/conv"
	gtest "github.com/og/x/test"
	"log"
	"net/http"
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
func getCookieCount(r *http.Request) (value int, err error) {
	c, err := r.Cookie("count")
	if err != nil {
		if err == http.ErrNoCookie {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return gconv.StringInt(c.Value)
}
func setCookieCount(w http.ResponseWriter, value int) {
	http.SetCookie(w,&http.Cookie{
		Name: "count",
		Value: gconv.IntString(value),
	})
}
func NewServe() *juice.Serve {
	serve := juice.NewServe(juice.ServeOption{
		OnCatchError: func(c *juice.Context, errInterface interface{}) error {
			log.Print(errInterface)
			return nil
		},
	})
	serve.HandleFunc(juice.POST, "/", func(c *juice.Context) (reject error) {
		req := ReqHome{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		reply := ReplyHome{}
		reply.IDNameAge  = req.ID + ":" + req.Name + ":" + gconv.IntString(req.Age)
		return c.Bytes(ogjson.Bytes(reply))
	})
	serve.HandleFunc(juice.POST, "/cookie", func(c *juice.Context) (reject error) {
		count , reject := getCookieCount(c.R) ; if reject != nil {return}
		newCount := count + 1
		setCookieCount(c.W, newCount)
		return c.Bytes([]byte("request:" + gconv.IntString(count) + " response:"+ gconv.IntString(newCount)))
	})
	return serve
}
func TestTest(t *testing.T) {
	as := gtest.NewAS(t)
	jtest := juicetest.NewTest(t, NewServe())
	resp := jtest.RequestJSON(juice.POST, "/", ReqHome{
		ID:   "a",
		Name: "b",
		Age:  3,
	})
	resp.ExpectJSON(200, ReplyHome{IDNameAge: "a:b:3",})
	{
		resp := jtest.RequestJSON(juice.POST, "/", ReqHome{
			ID: "a",
			Name: "b",
			Age: 3,
		})
		// resp.String resp.Bytes resp.BindJSON resp.ExpectJSON 任选其一即可
		as.Equal(resp.String(200), `{"idNameAge":"a:b:3"}`)

		reply := ReplyHome{}
		resp.BindJSON(200, &reply)
		as.Equal(reply, ReplyHome{
			IDNameAge: "a:b:3",
		})
		resp.ExpectJSON(200, ReplyHome{IDNameAge: "a:b:3",})
	}
	{
		jtest.RequestJSON(
			juice.POST, "/cookie", nil,
		).ExpectString(200, "request:0 response:1")
	}
	{
		jtest.RequestJSON(
			juice.POST, "/cookie", nil,
		).ExpectString(200, "request:1 response:2")
	}
}