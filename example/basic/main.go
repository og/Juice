package main

import (
	"github.com/michaeljs1990/sqlitestore"
	"github.com/og/juice"
	"log"
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

func main() {
	serve := juice.NewServe(juice.ServeOption{
		Session: sessionStore,
		OnCatchError: func(c *juice.Context, errInterface interface{}) error {
			log.Print(errInterface)
			switch errInterface.(type) {
			case error:
				err := errInterface.(error)
				return c.Bytes([]byte(err.Error()))
			default:
				return c.Bytes([]byte("server error"))
			}
		},
	})
	requestLogMiddleware := func(c *juice.Context, next juice.Next) error {
		log.Print(c.R.Method, " ", c.R.URL)
		return next()
	}
	serve.Use(requestLogMiddleware)
	serve.Action(juice.GET, "/", func(c *juice.Context) (reject error) {
		/* 绑定请求 */{
			req := struct {
				Name string `json:"name"`
				Age uint `json:"age"`
			}{}
			reject = c.BindRequest(&req) ;if reject != nil {return}
		}
		/* 读写 session */{
			sess := c.Session()
			// sess.SetString("time", time.Now().String())
			var timeStr string
			timeStr, reject = sess.GetString("time") ; if reject != nil {return}
			return c.Bytes([]byte("time:" + timeStr))
		}
	})
	{
		g := serve.Group()
		g.Use(func(c *juice.Context, next juice.Next) error {
			token := c.R.Header.Get("token")
			if token != "abc" {
				return c.Bytes([]byte("token 错误"))
			}
			return next()
		})
		g.Action(juice.GET, "/user", func(c *juice.Context) error {
			return c.Bytes([]byte("some data"))
		})
	}
	err := serve.Listen(":1219"); if err != nil {panic(err)}

}