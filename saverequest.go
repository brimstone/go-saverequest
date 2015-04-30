package saverequest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	RequestDir    = "."
	WriteRequests = false
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func Save(r *http.Request) {
	var body []byte
	if r.ContentLength > 0 {
		body, _ = ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
	r.Body = nopCloser{bytes.NewBufferString(string(body))}
	if !WriteRequests {
		return
	}
	log.Println("Saving", r.URL.Path)
	// create the proper directory structure for our request
	dir := RequestDir + r.URL.Path
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Println("Unable to make " + dir + ": " + err.Error())
		return
	}
	// create our file
	filename := strings.Replace(r.Header.Get("User-Agent"), "/", "_", -1)
	filename = dir + "/" + filename + "_request"
	output, err := os.Create(filename)
	if err != nil {
		log.Println("Unable to make " + filename + ": " + err.Error())
		return
	}

	// actually write our request
	fmt.Fprintf(output, "%s %s %s\n", r.Method, r.URL.Path, r.Proto)
	for k, v := range r.Header {
		for x := range v {
			fmt.Fprintf(output, "%s: %s\n", k, v[x])
		}
	}
	fmt.Fprintf(output, "\n")
	fmt.Fprintf(output, "%s", body)
	return
}
