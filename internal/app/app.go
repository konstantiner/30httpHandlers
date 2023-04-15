package app

import (
	"30httpHandlers/internal/entities"
	"30httpHandlers/internal/services"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() {
	rtr := chi.NewRouter()
	rtr.Use(middleware.Logger)

	rtr.Get("/users", getAll)
	rtr.Post("/users", create)
	rtr.Post("/friends", makeFriends)
	rtr.Delete("/users", deleteUser)
	rtr.Get("/users/{userID}/friends", userFriends)
	rtr.Put("/users/{userID}/age", updateUserAge)
	
	http.ListenAndServe("localhost:8080", rtr)
}

//GetAll возвращает всех пользователей в json формате
func getAll(w http.ResponseWriter, r *http.Request) {
	b := services.GetAllUsers()
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func create(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	var u entities.User
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	b := services.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func makeFriends(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	var u entities.MakeFriends
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	b := services.NewFriends(u.SourceId, u.TargetId)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	var u entities.MakeFriends
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	b := services.DeleteUser(u.TargetId)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func userFriends(w http.ResponseWriter, r *http.Request){
	userID, _:= strconv.Atoi(chi.URLParam(r, "userID"))
	
	b := services.UserFriends(userID)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func updateUserAge(w http.ResponseWriter, r *http.Request){
	type newAge struct{
		Age int `json:"age"`
	}

	userId, _ := strconv.Atoi(chi.URLParam(r, "userID"))
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	var u newAge
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	b := services.UpdateUserAge(userId, u.Age)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}