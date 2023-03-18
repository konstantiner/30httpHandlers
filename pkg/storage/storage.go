package storage

import (
	"30httpHandlers/pkg/users"
)

type Service struct {
	store map[int]users.User
}

func NewService() *Service {
	return &Service{
		store: map[int]users.User{},
	}
}

var counterId int = 1

func (s *Service) AllUsers() map[int]users.User{
	return s.store
}

func (s *Service) CreateUser(u users.User) int {
	newId := counterId
	counterId ++
	s.store[newId] = u
	return newId
}

func (s *Service) MakeFriends(SourceId int, TargetId int) {

	if entry, ok := s.store[SourceId]; ok{
		entry.Friends = append(entry.Friends, TargetId)
		s.store[SourceId] = entry
	}

	if entry, ok := s.store[TargetId]; ok{
		entry.Friends = append(entry.Friends, SourceId)
		s.store[TargetId] = entry
	}
}

func (s *Service) DeleteUser(userId int) {
	delete(s.store, userId) 
}

func (s *Service) DeleteFriend(userId int, friendId int) {
	var friendsUser []int = s.store[userId].Friends
	
	for y, z := range friendsUser {
		if friendId == z {
			friendsUser = append(friendsUser[:y], friendsUser[y+1:]... )
			if entry, ok := s.store[userId]; ok{
				entry.Friends = friendsUser
				s.store[userId] = entry
			}
			break
		}
	}
}

func (s *Service) UserFriends(userId int) []int {
	friends := s.store[userId].Friends
	return friends
}

func (s *Service) UserName(userId int) string {
	name := s.store[userId].Name
	return name
}

func (s *Service) UpdateUserAge(userId int, age int) {
	if entry, ok := s.store[userId]; ok{
		entry.Age = age
		s.store[userId] = entry
	}
}