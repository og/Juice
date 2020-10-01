package juice

import (
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type SessionStore interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
}

type Session struct {
	name string
	store SessionStore
	c *Context
}
func NewSession(name string, sessionStore SessionStore, c *Context) Session {
	return Session{
			name: name,
			store: sessionStore,
			c: c,
	}
}
func (s Session) GetString(key string) (value string, has bool, err error) {
	sess, err := s.store.Get(s.c.R, s.name)
	if err != nil {return }
	valueInterface := sess.Values[key]
	log.Print(valueInterface)
	switch valueInterface.(type) {
	case string:
		return valueInterface.(string), true, nil
	case nil:
		return "",false, nil
	default:
		return "", false, errors.New("juice.Session{}.GetString(key string)(value string err error) value type is not string")
	}
}
func (s Session) SetString(key string, value string) (err error) {
	sess, err := s.store.Get(s.c.R, s.name)
	sess.Values[key] = value
	err = sess.Save(s.c.R, s.c.W) ; if err != nil {return}
	return nil
}


func (c *Context) Session() Session {
	return NewSession("juice_session", c.serve.session, c)
}