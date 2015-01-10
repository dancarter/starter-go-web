package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/pat"
)

var logger = func(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}

func HomePageHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hello World!")
}

func Logger1Middleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	logger("[LOGGER 1] Started %s %s", req.Method, req.URL.Path)
	next(res, req)
	logger("[LOGGER 1] Completed in %v", time.Since(start))
}

func Logger2Middleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	logger("[LOGGER 2] Started %s %s", req.Method, req.URL.Path)
	next(res, req)
	logger("[LOGGER 2] Completed in %v", time.Since(start))
}

type MiddlewareHandler func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc)

type MeWare struct {
	Middles []http.HandlerFunc
}

func (self *MeWare) Use(middle MiddlewareHandler) {
	m := self.Middles[len(self.Middles)-1]
	sh := func(res http.ResponseWriter, req *http.Request) {
		middle(res, req, m)
	}
	self.Middles = append(self.Middles, sh)
}

func (self *MeWare) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	self.Middles[len(self.Middles)-1].ServeHTTP(res, req)
}

func NewMeWare(router http.Handler) *MeWare {
	m := &MeWare{}
	m.Middles = []http.HandlerFunc{router.ServeHTTP}
	return m
}

func Router() http.Handler {
	p := pat.New()
	p.Get("/", HomePageHandler)

	m := NewMeWare(p)
	m.Use(Logger1Middleware)
	m.Use(Logger2Middleware)
	return m
}

func main() {
	http.ListenAndServe(":3000", Router())
}
