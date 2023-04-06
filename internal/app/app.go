package app

import (
	"30httpHandlers/internal/entities"
	"30httpHandlers/internal/services"
	"encoding/json"
	"fmt"
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

	rtr.Get("/users", a.GetAll)
	rtr.Post("/users", a.Create)
	rtr.Post("/friends", a.MakeFriends)
	rtr.Delete("/users", a.DeleteUser)
	rtr.Get("/users/{userID}/friends", a.UserFriends)
	rtr.Put("/users/{userID}/age", a.UpdateUserAge)
	
	http.ListenAndServe("localhost:8080", rtr)
}

//GetAll возвращает всех пользователей в json формате
func (a *App) GetAll(w http.ResponseWriter, r *http.Request) {
	var response []entities.User
	for _, user := range a.repository.AllUsers() {
			response = append(response, user)
	}
	jsonString, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonString))
}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
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
	
	userId := a.repository.CreateUser(u)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Пользователь %s добавлен. ID = %d", u.Name, userId)))
}

func (a *App) MakeFriends(w http.ResponseWriter, r *http.Request) {
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
	rep := *&a.repository
	
	b:= services.NewFriends(u.SourceId, u.TargetId, rep)

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	//вытащим список друзей пользователя, зайдем к каждому другу и удалим у него удаленного пользователя из списка друзей
	var allFriends []int = a.repository.UserFriends(u.TargetId)
	for _, x := range allFriends {
		var friendsUser []int = a.repository.UserFriends(x)
		for _, z := range friendsUser {
			if u.TargetId == z {
				a.repository.DeleteFriend(x, z)				
			}
		}
	}
	nameDeleteUser :=  a.repository.UserName(u.TargetId)
	a.repository.DeleteUser(u.TargetId)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Пользователь %s удален.", nameDeleteUser)))
}

func (a *App) UserFriends(w http.ResponseWriter, r *http.Request){
	userID, _:= strconv.Atoi(chi.URLParam(r, "userID"))
	
	var allFriends []int = a.repository.UserFriends(userID)
	friendsName := ""
	for _, x := range allFriends {
		friendsName += fmt.Sprintf("ID: %d, Имя: %s\n",x, a.repository.UserName(x))
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Друзья пользователя: %s\n", friendsName)))
}

func (a *App) UpdateUserAge(w http.ResponseWriter, r *http.Request){
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
	a.repository.UpdateUserAge(userId, u.Age)
			
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Возраст пользователя успешно обновлён"))
}