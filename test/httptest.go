package juicetest

import (
	"bytes"
	ogjson "github.com/og/json"
	"github.com/og/juice"
	ge "github.com/og/x/error"
	gtest "github.com/og/x/test"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)


type Response struct {
	t *testing.T
	recorder *httptest.ResponseRecorder
	HttpResponse *http.Response
	as *gtest.AS
}
func (resp *Response) Bytes(statusCode int) []byte {
	resp.as.Equal(statusCode, resp.HttpResponse.StatusCode)
	b, err := ioutil.ReadAll(resp.recorder.Body) ; ge.Check(err)
	resp.recorder.Body = bytes.NewBuffer(b)
	return b
}
func (resp *Response) String(statusCode int) string {
	return string(resp.Bytes(statusCode))
}
func (resp *Response) ExpectString(statusCode int, s string) {
	resp.as.Equal(s, string(resp.Bytes(statusCode)))
}
func (resp *Response) BindJSON(statusCode int, v interface{})  {
	ogjson.ParseBytes(resp.Bytes(statusCode), v)
}
func (resp *Response) ExpectJSON(statusCode int, reply interface{}) {
	resp.as.Equal(ogjson.String(reply), resp.String(statusCode))

}
type Test struct {
	serve *juice.Serve
	t *testing.T
	jar *cookiejar.Jar
}
func NewTest(t *testing.T, serve *juice.Serve) Test {
	jar, err := cookiejar.New(nil) ; ge.Check(err)
	return Test{
		serve: serve,
		t: t,
		jar: jar,
	}
}
func (test Test) RequestJSON(method juice.Method, path string, jsonValue interface{}) (resp *Response)  {
	request := NewRequestJSON(method, path, jsonValue)
	return test.Request(request)
}
func (test *Test) Request(r *http.Request) (resp *Response)  {
	r.URL.Scheme = "http"
	/* request set cookie */{
		cookies := test.jar.Cookies(r.URL)
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
	}
	router := test.serve.HttpTestRouter()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, r)
	httpResponse :=  recorder.Result()
	/* response set cookie */ {
		test.jar.SetCookies(r.URL, httpResponse.Cookies())
	}
	return &Response{
		t: test.t,
		recorder: recorder,
		HttpResponse: httpResponse,
		as: gtest.NewAS(test.t),
	}
}

func NewRequestJSON(method juice.Method, path string, jsonValue interface{}) *http.Request {
	request := httptest.NewRequest(method.String(), path, bytes.NewReader(ogjson.Bytes(jsonValue)))
	request.Header.Set("Content-Type", "application/json")
	return request
}
