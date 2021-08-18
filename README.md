# Postier

Save and then serve all those POST requests received on a defined HTTP endpoint (WIP).

## Build

```shell
make
```

## Run

```shell
./bin/postier
2021/08/12 17:09:36 Warning, no LISTEN_URL provided, using default one (0.0.0.0:8042)
2021/08/12 17:09:36 Warning, no HISTORY_ENDPOINT provided, using default one (/postier-history)
2021/08/12 17:09:36 Info, history url /postier-history
```

## Usage

### As a standalone postier server

Make post requests

```shell
curl -d "test" -v http://localhost:8042/test
*   Trying 127.0.0.1:8042...
* Connected to localhost (127.0.0.1) port 8042 (#0)
> POST /test HTTP/1.1
> Host: localhost:8042
> User-Agent: curl/7.78.0
> Accept: */*
> Content-Length: 4
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 12 Aug 2021 15:14:48 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

Fetch history

```shell
curl http://localhost:8042/postier-history | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   254  100   254    0     0   315k      0 --:--:-- --:--:-- --:--:--  248k
[
  {
    "timestamp": "2021-08-16T15:40:59.32346424+02:00",
    "host": "localhost:8042",
    "url": "/test",
    "headers": {
      "Accept": [
        "*/*"
      ],
      "Content-Length": [
        "4"
      ],
      "Content-Type": [
        "application/x-www-form-urlencoded"
      ],
      "User-Agent": [
        "curl/7.78.0"
      ]
    },
    "method": "POST",
    "body": "test"
  }
]
```

### Configuration

- HISTORY_ENDPOINT environment variable is used to change the name of the main endpoint (in case you need to post data to this specific endpoint)
- LISTEN_URL environment variable is used to set the address and listening port (8042 is the default) of the server

### As a lib in Golang unit tests

Here are the basic for testing stuff using postier

1. start the server using `postester.StartTesting()` (this will start the server on a random available port), save the returned structure of kind `postester.Env`

2. post data to the running server (url is available in the `postester.Env` struct)

3. fetch history using `History()` method attached on the `postester.Env` struct

4. assert and search for you requests in history (you can pre-filter requests using the available methods)

See `examples/example_test.go`

```go
package examples

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/factorysh/postier/pkg/postester"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	// start the server
	pt, err := postester.StartTesting()
	assert.NoError(t, err)
	// cleanup at the end of the test
	defer pt.Cleanup()

	// fake post data
	values := map[string]string{"key": "value"}
	data, err := json.Marshal(values)
	assert.NoError(t, err)

	// post request with fake data
	resp, err := http.Post(fmt.Sprintf("http://%s/test", pt.URL), "application/json", bytes.NewReader(data))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// ask for posted data
	requests := pt.History().FilterURL("/test")
	assert.Len(t, requests, 1)
}

```
