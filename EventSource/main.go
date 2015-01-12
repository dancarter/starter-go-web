package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/antage/eventsource"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
)

const Admin = "Admin"

type UserList map[string]string

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (self Event) ToJSON() string {
	b, _ := json.Marshal(self)
	return string(b)

}

type Message struct {
	Name    string `json:"name"`
	Message string `json:"msg"`
}

func (self UserList) MarshalJSON() ([]byte, error) {
	names := []string{}
	for _, name := range self {
		names = append(names, name)
	}
	return json.Marshal(names)
}

func MessagesHandler(res http.ResponseWriter, req *http.Request) {
	msg := req.FormValue("msg")
	name := req.FormValue("name")
	sendMessage(name, msg)
}

func DeleteUsersHandler(res http.ResponseWriter, req *http.Request) {
	name := req.FormValue("user")
	if strings.ToLower(name) != "admin" {
		delete(Users, strings.ToLower(name))
		sendMessage(Admin, name+" has entered the room.")
		sendUsers()
	}
}

func PostUsersHandler(res http.ResponseWriter, req *http.Request) {
	name := req.FormValue("user")
	if strings.ToLower(name) != "admin" {
		Users[strings.ToLower(name)] = name
		sendMessage(Admin, name+" has entered the room.")
		sendUsers()
	}
}

func sendMessage(from string, msg string) {
	Streamer <- Event{Data: Message{Name: from, Message: msg}, Type: "message"}
}

func sendUsers() {
	Streamer <- Event{Data: Users, Type: "users"}
}

func processEvents(es eventsource.EventSource) {
	for {
		event := <-Streamer
		es.SendEventMessage(event.ToJSON(), "", strconv.Itoa(time.Now().Nanosecond()))
	}
}

var Streamer chan Event
var Users UserList

func main() {
	Streamer = make(chan Event)

	Users = make(map[string]string)
	Users["admin"] = Admin

	es := eventsource.New(nil, nil)
	defer es.Close()

	go processEvents(es)

	p := pat.New()
	p.Post("/messages", MessagesHandler)
	p.Post("/users", PostUsersHandler)
	p.Delete("/users", DeleteUsersHandler)
	p.Handle("/stream", es)

	n := negroni.Classic()
	n.UseHandler(p)

	http.ListenAndServe(":3000", n)
}
