package juice

import (
	"bytes"
	ogjson "github.com/og/json"
	ge "github.com/og/x/error"
	gtest "github.com/og/x/test"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func HttpTestRequestJSON(method Method, path string, v interface{}) (r *http.Request) {
	request, err := http.NewRequest(method.String(), path, bytes.NewReader(ogjson.Bytes(v))) ; ge.Check(err)
	request.Header.Set("Content-Type", "application/json")
	return  request
}

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

func HttpTest (addr string,r *http.Request) (resp HttpTestResponse)  {
	var err error
	r.URL, err = url.Parse("http://127.0.0.1" + addr + r.URL.Path) ; ge.Check(err)
	client := http.Client{}
	resp.HttpResponse, err = client.Do(r) ; ge.Check(err)
	return
}
