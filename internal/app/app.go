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

type App struct {
	repository services.Storage
}

func NewApp(repository services.Storage) *App {
	return &App{
		repository: repository,
	}
}

func (a *App) Run() {
	rtr := chi.NewRouter()
	rtr.Use(middleware.Logger)

	rtr.Get("/users", a.getAll)
	rtr.Post("/users", a.create)
	rtr.Post("/friends", a.makeFriends)
	rtr.Delete("/users", a.deleteUser)
	rtr.Get("/users/{userID}/friends", a.userFriends)
	rtr.Put("/users/{userID}/age", a.updateUserAge)
	
	http.ListenAndServe("localhost:8080", rtr)
}

//GetAll возвращает всех пользователей в json формате
func (a *App) getAll(w http.ResponseWriter, r *http.Request) {
	b := services.GetAllUsers(a.repository)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (a *App) create(w http.ResponseWriter, r *http.Request) {
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
	
	b := services.CreateUser(u, a.repository)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func (a *App) makeFriends(w http.ResponseWriter, r *http.Request) {
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
	
	b := services.NewFriends(u.SourceId, u.TargetId, a.repository)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
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

	b := services.DeleteUser(u.TargetId, a.repository)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (a *App) userFriends(w http.ResponseWriter, r *http.Request){
	userID, _:= strconv.Atoi(chi.URLParam(r, "userID"))
	
	b := services.UserFriends(userID, a.repository)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (a *App) updateUserAge(w http.ResponseWriter, r *http.Request){
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
	
	b := services.UpdateUserAge(userId, u.Age, a.repository)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}