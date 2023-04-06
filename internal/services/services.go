package services

import (
	"fmt"
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