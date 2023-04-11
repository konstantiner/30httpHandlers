package services

import (
	"fmt"
	"encoding/json"
	"30httpHandlers/internal/entities"
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

func GetAllUsers(rep Storage) (b []byte){
	var response []entities.User
	for _, user := range rep.AllUsers() {
			response = append(response, user)
	}
	b, _ = json.Marshal(response)
	return b
}

func CreateUser(u entities.User, rep Storage) (b []byte) {
	userId := rep.CreateUser(u)
	b = []byte(fmt.Sprintf("Пользователь %s добавлен. ID = %d", u.Name, userId))
	return
}

func NewFriends(SourceId int, TargetId int, rep Storage) (b []byte){
	username1 := rep.UserName(SourceId) 
    username2 := rep.UserName(TargetId) 
	
	//проверяем: если уже дружат, то не будем ещё раз добавлять в список друзей
    var allFriendsUser []int = rep.UserFriends(SourceId)
    for _, x := range allFriendsUser {
        if x == TargetId {
            b = []byte(fmt.Sprintf("%s и %s уже дружат", username1, username2))
            return
        }
    }

    rep.MakeFriends(SourceId, TargetId)
	b = []byte(fmt.Sprintf("%s и %s теперь друзья", username1, username2))
	return
}

func DeleteUser(TargetId int, rep Storage) (b []byte) {
	//вытащим список друзей пользователя, зайдем к каждому другу и удалим у него удаленного пользователя из списка друзей
	var allFriends []int = rep.UserFriends(TargetId)
	for _, x := range allFriends {
		var friendsUser []int = rep.UserFriends(x)
		for _, z := range friendsUser {
			if TargetId == z {
				rep.DeleteFriend(x, z)				
			}
		}
	}
	nameDeleteUser :=  rep.UserName(TargetId)
	rep.DeleteUser(TargetId)

	b = []byte(fmt.Sprintf("Пользователь %s удален.", nameDeleteUser))
	return
}

func UserFriends(userID int, rep Storage) (b []byte){
	var allFriends []int = rep.UserFriends(userID)
	friendsName := ""
	for _, x := range allFriends {
		friendsName += fmt.Sprintf("ID: %d, Имя: %s\n",x, rep.UserName(x))
	}

	b = []byte(fmt.Sprintf("Друзья пользователя: %s\n", friendsName))
	return
}

func UpdateUserAge(userId int, age int,rep Storage) (b []byte) {
	rep.UpdateUserAge(userId, age)
	b = []byte("Возраст пользователя успешно обновлён")
	return
}