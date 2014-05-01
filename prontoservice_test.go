package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	gc "gopkg.in/v1/check"
)

type ProntoServiceSuite struct{}

var _ = gc.Suite(&ProntoServiceSuite{})

type testSwift struct {
	err       error
	headers   http.Header
	container string
	object    string
	content   []byte
}

func NewTestSwift(h http.Header, err error, c []byte) *testSwift {
	return &testSwift{err: err, headers: h, content: c}
}

func (t *testSwift) GetReader(container, object string) (io.ReadCloser, http.Header, error) {
	t.container = container
	t.object = object
	return &swiftFile{Buffer: bytes.NewBuffer(t.content)}, t.headers, t.err
}

type swiftFile struct {
	*bytes.Buffer
	closed bool
}

func (s *swiftFile) Close() error {
	s.closed = true
	return nil
}

// ProntoService should return a 404 if there's an error from Swift service
func (s *ProntoServiceSuite) TestReturns404WhenErrorReceivedFromSwift(c *gc.C) {
	service := ProntoService{NewTestSwift(nil, errors.New("Test Error"), nil), "testing"}

	req, err := http.NewRequest("GET", "/my-test-url", nil)
	c.Assert(err, gc.IsNil)

	w := httptest.NewRecorder()
	service.ServeHTTP(w, req)

	c.Assert(w.Code, gc.Equals, 404)
}

// ProntoService passes through the returned headers from the Swift service
func (s *ProntoServiceSuite) TestProxiesHeadersFromSwift(c *gc.C) {
	headers := http.Header{"Etag": []string{"\"testing\""}}
	service := ProntoService{NewTestSwift(headers, nil, nil), "testing"}

	req, err := http.NewRequest("GET", "/my-test-url", nil)
	c.Assert(err, gc.IsNil)

	w := httptest.NewRecorder()
	service.ServeHTTP(w, req)

	c.Assert(w.Header(), gc.DeepEquals, headers)
}

// ProntoService chops the initial / from the URL before passing it to Swift
func (s *ProntoServiceSuite) TestPathPassedToSwift(c *gc.C) {
	testSwift := NewTestSwift(nil, nil, nil)
	service := ProntoService{testSwift, "testing"}
	req, err := http.NewRequest("GET", "/my-test-url", nil)
	c.Assert(err, gc.IsNil)

	w := httptest.NewRecorder()
	service.ServeHTTP(w, req)

	c.Assert(testSwift.container, gc.Equals, "testing")
	c.Assert(testSwift.object, gc.Equals, "my-test-url")
}

// ProntoService writes the content of the Swift file to the response
func (s *ProntoServiceSuite) TestSwiftContentSentToBrowser(c *gc.C) {
	testSwift := NewTestSwift(nil, nil, []byte(`This is a test`))
	service := ProntoService{testSwift, "testing"}
	req, err := http.NewRequest("GET", "/my-test-url", nil)
	c.Assert(err, gc.IsNil)

	w := httptest.NewRecorder()
	service.ServeHTTP(w, req)

	c.Assert(w.Body.String(), gc.Equals, "This is a test")
}

// TODO: provide a log to ProntoService and test New()
