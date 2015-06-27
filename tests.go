package saverequest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func FakeRequest(method string, path string, headers map[string]string, body string) (*http.Request, error) {
	// build our request
	req, err := http.NewRequest(method, "http://example.com"+path, nil)
	req.RequestURI = path
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// TODO add in body
	req.Body = nopCloser{bytes.NewBufferString(body)}
	req.ContentLength = int64(len(body))
	// return our request
	return req, nil
}

func newRequestFromFile(filename string) (*http.Request, error) {
	// read in our file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// split it up and grab the method, and path
	lines := strings.Split(string(content), "\n")
	firstline := strings.Split(lines[0], " ")
	method := firstline[0]
	path := firstline[1]

	// add in headers
	headers := make(map[string]string)
	var l = 1
	for ; lines[l] != ""; l++ {
		header := strings.SplitN(lines[l], ": ", 2)
		headers[header[0]] = header[1]
	}

	body := content[len(strings.Join(lines[:l+1], "\n"))+1:]
	return FakeRequest(method, path, headers, string(body))

}

func TestRequestFiles(t *testing.T, basedir string, handler func(http.ResponseWriter, *http.Request)) {
	var files []string

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// if we have a request
		if strings.HasSuffix(path, "_request") {
			request := strings.TrimSuffix(path, "_request")
			response := request + "_response"
			file, err := os.Open(response)
			// and a matching response
			if err == nil {
				// save it for later
				files = append(files, request)
				file.Close()
			}
		}
		return nil
	})

	// now actually check our requests with our responses
	for f := range files {
		fmt.Printf("File: %s\n", files[f])
		req, err := newRequestFromFile(files[f] + "_request")
		if err != nil {
			t.Errorf(err.Error())
		}
		response, err := ioutil.ReadFile(files[f] + "_response")
		if err != nil {
			t.Errorf(err.Error())
		}
		w := httptest.NewRecorder()
		handler(w, req)
		if w.Code != 200 {
			t.Errorf("Response is not 200")
		}

		if w.Body.String() != string(response) {
			t.Errorf("Response is not identical\nExpected: \"%s\"\nReceived: \"%s\"\n", string(response), w.Body.String())
		}
		fmt.Printf("%d - %s\n", w.Code, w.Body.String())
	}
}
