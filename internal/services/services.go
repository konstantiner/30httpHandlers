package services

import (
	"fmt"
)

func NewFriends(SourceId int, TargetId int, username1 string, username2 string, allFriendsUser1 []int) (b []byte, chek bool){
	
	//проверяем: если уже дружат, то не будем ещё раз добавлять в список друзей
	for _, x := range allFriendsUser1 {
		if x == TargetId {
			b = []byte(fmt.Sprintf("%s и %s уже дружат", username1, username2))
			return b, true
		}
	}
	
	b = []byte(fmt.Sprintf("%s и %s теперь друзья", username1, username2))
	return b, false
}