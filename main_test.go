package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	client = &http.Client{Timeout: time.Second}
)

type Case struct {
	Method string // GET по-умолчанию в http.NewRequest если передали пустую строку
	Path   string
	Query  string
	Status int
	// Result interface{}
}

//CaseResponce
type CR map[string]interface{}

func TestServer(t *testing.T) {
	cases := []*Case{
		&Case{
			"", AllPreparationsURL, "", http.StatusOK,
		},
	}
	serv := httptest.NewServer(NewHandler())

	for idx, item := range cases {
		var (
			err error
			// result   interface{}
			// expected interface{}
			req *http.Request
		)

		caseName := fmt.Sprintf("case %d: [%s] %s %s", idx, item.Method, item.Path, item.Query)

		if item.Method == http.MethodPost {
			reqBody := strings.NewReader(item.Query)
			req, err = http.NewRequest(item.Method, serv.URL+item.Path, reqBody)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req, err = http.NewRequest(item.Method, serv.URL+item.Path+"?"+item.Query, nil)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("[%s] request error: %v", caseName, err)
			continue
		}
		// defer resp.Body.Close()

		if resp.StatusCode != item.Status {
			t.Errorf("[%s] expected http status %v, got %v", caseName, item.Status, resp.StatusCode)
			continue
		}
	}
}
