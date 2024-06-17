package main

import (
	"sync"
	"time"
)

type profiles struct {
	mu sync.RWMutex
	m  map[string]UserProfile
}

type UserProfile struct {
	UUID   string
	Name   string
	Orders []Order
	TTL    time.Time
}

type Order struct {
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Basket    interface{}
}

func (p *profiles) Get(key string) (UserProfile, bool) {
	p.mu.RLock()
	v, ok := p.m[key]
	p.mu.RUnlock()
	if ok {
		return v, true
	}
	return UserProfile{}, false // или вернуть значение по умо лчанию
}

func (p *profiles) Add(user UserProfile) bool {
	if p.m == nil { // Проверяем, что карта является nil
		p.m = make(map[string]UserProfile) // Инициализируем карту
	}
	_, ok := p.m[user.UUID]
	if !ok {
		p.mu.Lock()
		user.TTL = time.Now().Add(time.Second * 3)
		p.m[user.UUID] = user
		p.mu.Unlock()

		return true
	}
	return false
}

func startTimer(p *profiles) {
	time.Sleep(time.Second * 5)
	t := time.NewTicker(15 * time.Second)
	for range t.C {
		for _, v := range p.m {
			if v.TTL.Before(time.Now()) {
				p.mu.Lock()
				delete(p.m, v.UUID)
				p.mu.Unlock()
			}
		}
	}
}
