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
	defer s.mutex.Unlock()
	newId := s.counterId
	s.counterId ++
	s.store[newId] = u
	return newId
}

func (s *Service) MakeFriends(SourceId int, TargetId int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.store, userId) 
}

func (s *Service) DeleteFriend(userId int, friendId int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	friends := s.store[userId].Friends
	return friends
}

func (s *Service) UserName(userId int) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	name := s.store[userId].Name
	return name
}

func (s *Service) UpdateUserAge(userId int, age int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if entry, ok := s.store[userId]; ok{
		entry.Age = age
		s.store[userId] = entry
	}
}