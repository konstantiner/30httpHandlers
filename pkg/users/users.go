package users

import (
	"fmt"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Friends []int `json:"friends"`
}

type MakeFriends struct {
	SourceId int `json:"sourceId"`
	TargetId int `json:"targetId"`
}

func ToString(id int, u User) string {
	return fmt.Sprintf("Id: %d, Name: %s, age: %d, friends: %d \n",id, u.Name, u.Age, u.Friends)
}