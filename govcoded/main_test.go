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

func TestRepos(t *testing.T) {
	res, _ := GetPath("/repos")

	if res.Code != 200 {
		t.Error("/repos should return 200")
	}
}

func TestOrgs(t *testing.T) {
	res, _ := GetPath("/orgs")

	if res.Code != 200 {
		t.Error("/orgs should return 200")
	}
}
