package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetPath(path string) (*httptest.ResponseRecorder, *http.Request) {
	m := App()
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	m.ServeHTTP(res, req)

	return res, req
}

func TestEndpoints(t *testing.T) {
	urls := []string{"/", "/repos", "/orgs", "/stats", "/issues"}

	for _, u := range urls {
		res, _ := GetPath(u)

		if res.Code != 200 {
			t.Error(u, "should return 200")
		}

	}
}
