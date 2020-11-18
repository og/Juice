package jhttptest_test

import (
	"fmt"
	"github.com/michaeljs1990/sqlitestore"
	ogjson "github.com/og/json"
	jhttp "github.com/og/juice/http"
	jhttptest "github.com/og/juice/http/test"
	gconv "github.com/og/x/conv"
	gtest "github.com/og/x/test"
	"log"
	"net/http"
	"testing"
)

var sessionStore *sqlitestore.SqliteStore
func init() {
	var err error
	sessionStore, err = sqlitestore.NewSqliteStore(
		"./test/session_sqllite",
		"sessions",
		"/",
		3600*24,
		[]byte("production environment must write secretKey"),
	)
	if err != nil {
		panic(err)
	}
}

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
func NewRouter() *jhttp.Router {
	router := jhttp.NewRouter(jhttp.RouterOption{
		OnCatchError: func(c *jhttp.Context, errInterface interface{}) error {
			log.Print(errInterface)
			return nil
		},
	})
	router.HandleFunc(jhttp.POST, "/", func(c *jhttp.Context) (reject error) {
		req := ReqHome{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		reply := ReplyHome{}
		reply.IDNameAge  = req.ID + ":" + req.Name + ":" + gconv.IntString(req.Age)
		return c.Bytes(ogjson.Bytes(reply))
	})
	router.HandleFunc(jhttp.POST, "/cookie", func(c *jhttp.Context) (reject error) {
		count , reject := getCookieCount(c.R) ; if reject != nil {return}
		newCount := count + 1
		setCookieCount(c.W, newCount)
		return c.Bytes([]byte("request:" + gconv.IntString(count) + " response:"+ gconv.IntString(newCount)))
	})
	router.HandleFunc(jhttp.GET, "/session_set", func(c *jhttp.Context) (reject error) {
		return c.Session("juice_session", sessionStore).SetString("userID", "a")
	})
	router.HandleFunc(jhttp.GET, "/session_get", func(c *jhttp.Context) (reject error) {
		userID, hasUserID, reject := c.Session("juice_session", sessionStore).GetString("userID") ; if reject != nil {return}
		return c.Bytes([]byte(fmt.Sprintf("%s,%v",userID,hasUserID)))
	})
	router.HandleFunc(jhttp.GET, "/session_del", func(c *jhttp.Context) (reject error) {
		return c.Session("juice_session", sessionStore).DelString("userID")
	})
	return router
}
func TestTest(t *testing.T) {
	as := gtest.NewAS(t)
	jtest := jhttptest.NewTest(t, NewRouter())
	resp := jtest.RequestJSON(jhttp.POST, "/", ReqHome{
		ID:   "a",
		Name: "b",
		Age:  3,
	})
	resp.ExpectJSON(200, ReplyHome{IDNameAge: "a:b:3",})
	{
		resp := jtest.RequestJSON(jhttp.POST, "/", ReqHome{
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
			jhttp.POST, "/cookie", nil,
		).ExpectString(200, "request:0 response:1")
	}
	{
		jtest.RequestJSON(
			jhttp.POST, "/cookie", nil,
		).ExpectString(200, "request:1 response:2")
	}
	{
		resp := jtest.RequestJSON(jhttp.GET, "/session_get", nil)
		resp.ExpectString(200, ",false")
	}
	{
		resp := jtest.RequestJSON(jhttp.GET, "/session_set", nil)
		resp.ExpectString(200,"")
	}
	{
		resp := jtest.RequestJSON(jhttp.GET, "/session_get", nil)
		resp.ExpectString(200, "a,true")
	}
	{
		resp := jtest.RequestJSON(jhttp.GET, "/session_del", nil)
		resp.ExpectString(200,"")
	}
	{
		resp := jtest.RequestJSON(jhttp.GET, "/session_get", nil)
		resp.ExpectString(200, ",false")
	}
}