package juice

import (
	"bytes"
	ogjson "github.com/og/json"
	ge "github.com/og/x/error"
	gtest "github.com/og/x/test"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)


type HttpTestResponse struct {
	HttpResponse *http.Response
}
func (resp *HttpTestResponse) Bytes() []byte {
	b, err := ioutil.ReadAll(resp.HttpResponse.Body) ; ge.Check(err)
	resp.HttpResponse.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return b
}
func (resp *HttpTestResponse) String() string {
	return string(resp.Bytes())
}
func (resp *HttpTestResponse) BindJSON(v interface{})  {
	ogjson.ParseBytes(resp.Bytes(), v)
}
func (resp *HttpTestResponse) ExpectJSON(t *testing.T, reply interface{}) {
	gtest.NewAS(t).Equal(ogjson.String(reply), resp.String())

}
type HttpTest struct {
	addr string
	client *http.Client
}
func NewHttpTest(addr string) HttpTest {
	jar, err := cookiejar.New(nil) ; ge.Check(err)
	return HttpTest{
		addr: addr,
		client: &http.Client{
			Jar: jar,
		},
	}
}
func (h HttpTest) RequestJSON(method Method, path string, jsonValue interface{}) (resp HttpTestResponse)  {
	var err error
	request, err := http.NewRequest(method.String(), path, bytes.NewReader(ogjson.Bytes(jsonValue))) ; ge.Check(err)
	request.Header.Set("Content-Type", "application/json")
	return h.Request(request)

}
func (h HttpTest) Request(r *http.Request) (resp HttpTestResponse)  {
	var err error
	r.URL, err = url.Parse("http://127.0.0.1" + h.addr + r.URL.Path) ; ge.Check(err)
	resp.HttpResponse, err = h.client.Do(r) ; ge.Check(err)
	return
}
