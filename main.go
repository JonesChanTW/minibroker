package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
)

type userObj struct {
	ID      int64
	Name    string
	Age     int
	Address string
	EMail   string
}

var userList = map[int64]*userObj{}
var gl = sync.RWMutex{}

func main() {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/users", getUserList),
		rest.Get("/user/:id", getUser),
		rest.Post("/CreateUser", createUser),
	)

	if err != nil {
		fmt.Println("Start App Error : ", err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":80", api.MakeHandler()))
}

func createUser(w rest.ResponseWriter, r *rest.Request) {
	user := userObj{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.ID != 0 {
		rest.Error(w, "Create Account can't identify id", 400)
		return
	}
	gl.Lock()
	user.ID = int64(len(userList) + 1)
	userList[user.ID] = &user
	gl.Unlock()
	w.WriteJson(user)
}

func getUserList(w rest.ResponseWriter, r *rest.Request) {
	rest.NotFound(w, r)
	return
}

func getUser(w rest.ResponseWriter, r *rest.Request) {
	strID := r.PathParam("id")
	id, err := strconv.ParseInt(strID, 10, 64)

	if err != nil {
		fmt.Println("getUser fail id = ", strID)
		rest.NotFound(w, r)
		return
	}

	gl.Lock()
	var user *userObj
	if userList[id] != nil {
		user = &userObj{}
		*user = *userList[id]
	}
	gl.Unlock()
	if user == nil {
		rest.NotFound(w, r)
		return
	}

	w.WriteJson(user)
}
