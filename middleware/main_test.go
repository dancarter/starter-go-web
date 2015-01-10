package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_App(t *testing.T) {
	a := assert.New(t)

	ts := httptest.NewServer(Router())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	a.NoError(err)

	body, err := ioutil.ReadAll(res.Body)
	a.NoError(err)
	a.Equal(string(body), "Hello World!\n")
}

func Test_Middleware(t *testing.T) {
	a := assert.New(t)

	logs := []string{}
	logger = func(msg string, args ...interface{}) {
		logs = append(logs, fmt.Sprintf(msg, args...))
	}

	ts := httptest.NewServer(Router())
	defer ts.Close()

	_, err := http.Get(ts.URL)
	a.NoError(err)

	a.NotEmpty(logs)
	a.Equal(len(logs), 4)

	// a.Contains(logs[0], "[LOGGER 2] Started")
	// a.Contains(logs[1], "[LOGGER 1] Started")
	// a.Contains(logs[2], "[LOGGER 1] Completed")
	// a.Contains(logs[3], "[LOGGER 2] Completed")
}
