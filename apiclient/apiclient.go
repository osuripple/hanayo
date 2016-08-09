// Package apiclient implements a minimal API client for the Ripple API,
// designed for hanayo's templates. Do not actually use this outside of Go's
// HTML templates.
package apiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// UserAgent is the User-Agent of the HTTP request.
const UserAgent = "hanayo"

// Key is H-Key, that will be passed to the API to tell not to rate limit this
// request.
var Key = "Potato"

// APIBase is the base url of the API.
var APIBase = "http://localhost:40001/api/v1/"

// Get retrieves data in a GET request to the Ripple API.
func Get(s string, params ...interface{}) (map[string]interface{}, error) {
	s = fmt.Sprintf(s, params...)
	req, err := http.NewRequest("GET", APIBase+s, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("H-Key", Key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	x := make(map[string]interface{})
	err = json.Unmarshal(data, &x)
	return x, err
}
