package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Friends []int `json:"friends"`
}

type makeFriends struct {
	SourceId int `json:"sourceId"`
	TargetId int `json:"targetId"`
}

type service struct {
	store map[int]*User
}

var counterId int = 1

func (u *User) toString(id int) string {
	return fmt.Sprintf("Id: %d, Name: %s, age: %d, friends: %d \n",id, u.Name, u.Age, u.Friends)
}

func main() {
	
	rtr := chi.NewRouter()
	rtr.Use(middleware.Logger)
	srv := service{make(map[int]*User)}

	rtr.Get("/get", srv.GetAll)
	rtr.Post("/create", srv.Create)
	rtr.Post("/make_friends", srv.MakeFriends)
	rtr.Delete("/delete", srv.DeleteUser)
	rtr.Get("/get/{userID}", srv.UserFriends)
	rtr.Put("/{userID}", srv.UpdateUserAge)
	
	http.ListenAndServe("localhost:8080", rtr)
}

func (s *service) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u User
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		s.store[counterId] = &u
		
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("Пользователь %s добавлен. ID = %d", u.Name, counterId)))
		counterId ++
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := ""
		for id, user := range s.store {
			response += user.toString(id)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u makeFriends
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		username1 := s.store[u.SourceId].Name
		username2 := s.store[u.TargetId].Name
		
		var allFriendsUser []int = s.store[u.SourceId].Friends
		for _, x := range allFriendsUser {
			if x == u.TargetId {
				w.WriteHeader(http.StatusAlreadyReported)
				w.Write([]byte(fmt.Sprintf("%s и %s уже дружат", username1, username2)))
				return
			}
		}
		
		s.store[u.SourceId].Friends = append(s.store[u.SourceId].Friends, u.TargetId) 
		s.store[u.TargetId].Friends = append(s.store[u.TargetId].Friends, u.SourceId)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s и %s теперь друзья", username1, username2)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u makeFriends
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

        //вытащим список друзей пользователя, зайдем к каждому другу и удалим у него удаленного пользователя из списка друзей
		var allFriends []int = s.store[u.TargetId].Friends
		for _, x := range allFriends {
			var friendsUser []int = s.store[x].Friends
			for y, z := range friendsUser {
				if u.TargetId == z {
					friendsUser = append(friendsUser[:y], friendsUser[y+1:]... )
					s.store[x].Friends = friendsUser
				}
			}
		}
		nameDeleteUser := s.store[u.TargetId].Name
			
		delete(s.store, u.TargetId) 
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Пользователь %s удален.", nameDeleteUser)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) UserFriends(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		userID, _:= strconv.Atoi(chi.URLParam(r, "userID"))
		
		var allFriends []int = s.store[userID].Friends
		friendsName := []string{}
		for _, x := range allFriends {
			friendsName = append(friendsName, s.store[x].Name)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Друзья пользователя: %s", friendsName)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) UpdateUserAge(w http.ResponseWriter, r *http.Request){
	if r.Method == "PUT" {
		type newAge struct{
			Age int `json:"age"`
		}

		userID, _ := strconv.Atoi(chi.URLParam(r, "userID"))
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

		s.store[userID].Age = u.Age
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Возраст пользователя успешно обновлён"))
	}
}

/*
curl -X POST -d "{\"name\": \"Vasiliy\", \"age\": 20}" http://localhost:8080/create
curl -X POST -d "{\"name\": \"Ivan\", \"age\": 30}" http://localhost:8080/create
curl -X POST -d "{\"name\": \"Boss\", \"age\": 30}" http://localhost:8080/create
curl -X POST -d "{\"sourceId\": 1, \"targetId\": 2}" http://localhost:8080/make_friends
curl -X DELETE -d "{\"targetId\": 1}" http://localhost:8080/delete
curl -X PUT -d "{\"age\": 25}" http://localhost:8080/2

http://localhost:8080/get
http://localhost:8080/get/1
*/