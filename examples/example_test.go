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
	resp, err := http.Post(fmt.Sprintf("%s/test", pt.URL), "application/json", bytes.NewReader(data))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// ask for posted data
	requests := pt.History().FilterURL("/test")
	assert.Len(t, requests, 1)
}
