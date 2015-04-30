go-saverequest
==============

This is a simple library that records requests on demand, then plays them back later during test time.

Usage
-----

Take for example, the following simple program:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/brimstone/go-saverequest"
)

func handleData(w http.ResponseWriter, r *http.Request) {

	saverequest.Save(r)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, "{\"message\": \"thanks!\"}")
}

func main() {

	saverequest.WriteRequests = true
	http.HandleFunc("/data/", handleData)
	http.ListenAndServe(":8000", nil)
}
```

... and its ever so simple test harness
```go
package main

import (
	"testing"

	"github.com/brimstone/go-saverequest"
)

func TestHandleData(t *testing.T) {
	saverequest.TestRequestFiles(t, ".", handleData)
}
```

This results in the following go test output:

```bash
$ go test
PASS
ok      example     0.006s
```

Now, the real magic happens by first running the program:

```bash
$ go run main.go
```

Then hitting the site with something like curl:

```bash
$ curl http://localhost:8000/data/file
{"message": "thanks!"}
```

You should see that the program logged the request:

```bash
$ go run main.go
2015/04/29 23:36:45 Saving /data/test to ./data/test
```

Specifically, there is now a request file in `data/test/curl_7.42.0_request`

For sake of a simple test, run the request again, this time saving the response to where we expect:

```bash
$ curl http://localhost:8000/data/file > data/test/curl_7.42.0_response
```

Now when you run the tests you should see something different:
```bash
File: data/test/curl_7.42.0
200 - {"message": "thanks!"}
PASS
ok      example     0.006s
```


Repeat as needed, season to taste.
