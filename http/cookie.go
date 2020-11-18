package jhttp

import (
	"net/http"
)

func (c *Context) Cookie() Cookie {
	return NewCookie(c.R, c.W)
}
type Cookie struct {
	r *http.Request
	w http.ResponseWriter
}
func NewCookie(r *http.Request, w http.ResponseWriter) Cookie {
	return Cookie{
		r: r,
		w: w,
	}
}
func (c Cookie) GetString(name string) (value string, has bool, reject error) {
	cookie, err := c.r.Cookie(name)
	switch err {
	case nil:
		return cookie.Value, true, nil
	case http.ErrNoCookie:
		return "", false, nil
	default:
		return "", false, err
	}
}
func (c Cookie) SetString(name string, value string) {
	http.SetCookie(c.w, &http.Cookie{
		Name:  name,
		Value: value,
	})
}
