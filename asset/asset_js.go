// +build js

package asset

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func loadFile(path string) ([]byte, error) {
	return httpGet(path)
}

// httpGet fetches the contents at the given url. This is used as a workaround for loading local assets while on the web.
// TODO: This prevents `gopherjs build` and then running locally because you can't use http.GET with local files. It has to be with http or https targets.
// I don't think there's a good way around this. Just don't use `gopherjs build`. Instead use `gopherjs serve` and view the served version.
// Consider instead returning an io.ReadCloser so large responses are better dealt with.
func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status: %s", resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
