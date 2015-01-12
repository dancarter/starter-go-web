package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/unrolled/render"
)

type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email_address" xml:"emailAddress"`
}

func main() {
	r := render.New(render.Options{
		IndentJSON:    true,
		Directory:     "templates",
		Extensions:    []string{".html", ".tmpl"},
		Layout:        "application",
		Funcs:         []template.FuncMap{{"upper": strings.ToUpper}},
		IndentXML:     true,
		PrefixXML:     []byte("<?xml version='1.0' encoding='UTF-8'?>\n'"),
		IsDevelopment: true,
	})

	http.HandleFunc("/json", func(res http.ResponseWriter, req *http.Request) {
		user := userFromReq(req)
		r.JSON(res, 200, user)
	})
	http.HandleFunc("/xml", func(res http.ResponseWriter, req *http.Request) {
		user := userFromReq(req)
		r.XML(res, 200, user)
	})
	http.HandleFunc("/admin", func(res http.ResponseWriter, req *http.Request) {
		user := userFromReq(req)
		r.HTML(res, 200, "admin/index", user)
	})
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		user := userFromReq(req)
		r.HTML(res, 200, "index", user)
	})

	http.ListenAndServe(":4000", nil)
}

func userFromReq(req *http.Request) *User {
	user := &User{
		Name: nameFromReq(req),
	}
	user.Email = fmt.Sprintf("%s@example.com", user.Name)
	return user
}

func nameFromReq(req *http.Request) string {
	name := req.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	return name
}
