package services

import (
	"fmt"
	"encoding/json"
	"30httpHandlers/internal/entities"
	"30httpHandlers/internal/storage"
)

type Storage interface {
	AllUsers() map[int]entities.User
	CreateUser(entities.User) int
	MakeFriends(int, int)
	DeleteUser(userId int)
	DeleteFriend(userId int, friendId int)
	UserFriends(userId int) []int
	UserName(userId int) string
	UpdateUserAge(userId int, age int)
}

type App struct {
	repository Storage
}

func NewApp(repository Storage) *App {
	return &App{
		repository: repository,
	}
}

var repository = storage.NewMemStore()
var ser = NewApp(repository)

func GetAllUsers() (b []byte){
	var response []entities.User
	for _, user := range Storage.AllUsers(ser.repository) {
			response = append(response, user)
	}
	b, _ = json.Marshal(response)
	return b
}

func CreateUser(u entities.User) (b []byte) {
	userId := ser.repository.CreateUser(u)
	b = []byte(fmt.Sprintf("Пользователь %s добавлен. ID = %d", u.Name, userId))
	return
}

func NewFriends(SourceId int, TargetId int) (b []byte){
	username1 := ser.repository.UserName(SourceId) 
    username2 := ser.repository.UserName(TargetId) 
	
	//проверяем: если уже дружат, то не будем ещё раз добавлять в список друзей
    var allFriendsUser []int = ser.repository.UserFriends(SourceId)
    for _, x := range allFriendsUser {
        if x == TargetId {
            b = []byte(fmt.Sprintf("%s и %s уже дружат", username1, username2))
            return
        }
    }

    ser.repository.MakeFriends(SourceId, TargetId)
	b = []byte(fmt.Sprintf("%s и %s теперь друзья", username1, username2))
	return
}

func DeleteUser(TargetId int) (b []byte) {
	//вытащим список друзей пользователя, зайдем к каждому другу и удалим у него удаленного пользователя из списка друзей
	var allFriends []int = ser.repository.UserFriends(TargetId)
	for _, x := range allFriends {
		var friendsUser []int = ser.repository.UserFriends(x)
		for _, z := range friendsUser {
			if TargetId == z {
				ser.repository.DeleteFriend(x, z)				
			}
		}
	}
	nameDeleteUser :=  ser.repository.UserName(TargetId)
	ser.repository.DeleteUser(TargetId)

	b = []byte(fmt.Sprintf("Пользователь %s удален.", nameDeleteUser))
	return
}

func UserFriends(userID int) (b []byte){
	var allFriends []int = ser.repository.UserFriends(userID)
	friendsName := ""
	for _, x := range allFriends {
		friendsName += fmt.Sprintf("ID: %d, Имя: %s\n",x, ser.repository.UserName(x))
	}

	b = []byte(fmt.Sprintf("Друзья пользователя: %s\n", friendsName))
	return
}

func UpdateUserAge(userId int, age int) (b []byte) {
	ser.repository.UpdateUserAge(userId, age)
	b = []byte("Возраст пользователя успешно обновлён")
	return
}