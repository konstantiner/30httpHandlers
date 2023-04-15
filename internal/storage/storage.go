package storage

import (
	"30httpHandlers/internal/entities"
	"sync"
)

type MemStorage struct {
	counterId int 
	store map[int]entities.User
	mutex sync.Mutex 
}

func NewMemStore() *MemStorage {
	return &MemStorage{
		store: map[int]entities.User{},
	}
}

func (s *MemStorage) AllUsers() map[int]entities.User{
	return s.store
}

func (s *MemStorage) CreateUser(u entities.User) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	newId := s.counterId
	s.counterId ++
	s.store[newId] = u
	return newId
}

func (s *MemStorage) MakeFriends(SourceId int, TargetId int) {
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

func (s *MemStorage) DeleteUser(userId int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.store, userId) 
}

func (s *MemStorage) DeleteFriend(userId int, friendId int) {
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

func (s *MemStorage) UserFriends(userId int) []int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	friends := s.store[userId].Friends
	return friends
}

func (s *MemStorage) UserName(userId int) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	name := s.store[userId].Name
	return name
}

func (s *MemStorage) UpdateUserAge(userId int, age int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if entry, ok := s.store[userId]; ok{
		entry.Age = age
		s.store[userId] = entry
	}
}