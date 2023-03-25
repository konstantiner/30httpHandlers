package storage

import (
	"30httpHandlers/internal/entities"
	"sync"
)

type Service struct {
	counterId int 
	store map[int]entities.User
	mutex sync.Mutex 
}

func NewService() *Service {
	return &Service{
		store: map[int]entities.User{},
	}
}

func (s *Service) AllUsers() map[int]entities.User{
	return s.store
}

func (s *Service) CreateUser(u entities.User) int {
	s.mutex.Lock()
	newId := s.counterId
	s.counterId ++
	s.store[newId] = u
	s.mutex.Unlock()
	return newId
}

func (s *Service) MakeFriends(SourceId int, TargetId int) {
	s.mutex.Lock()
	if entry, ok := s.store[SourceId]; ok{
		entry.Friends = append(entry.Friends, TargetId)
		s.store[SourceId] = entry
	}

	if entry, ok := s.store[TargetId]; ok{
		entry.Friends = append(entry.Friends, SourceId)
		s.store[TargetId] = entry
	}
	s.mutex.Unlock()
}

func (s *Service) DeleteUser(userId int) {
	s.mutex.Lock()
	delete(s.store, userId) 
	s.mutex.Unlock()
}

func (s *Service) DeleteFriend(userId int, friendId int) {
	s.mutex.Lock()
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
	s.mutex.Unlock()
}

func (s *Service) UserFriends(userId int) []int {
	s.mutex.Lock()
	friends := s.store[userId].Friends
	s.mutex.Unlock()
	return friends
}

func (s *Service) UserName(userId int) string {
	s.mutex.Lock()
	name := s.store[userId].Name
	s.mutex.Unlock()
	return name
}

func (s *Service) UpdateUserAge(userId int, age int) {
	s.mutex.Lock()
	if entry, ok := s.store[userId]; ok{
		entry.Age = age
		s.store[userId] = entry
	}
	s.mutex.Unlock()
}