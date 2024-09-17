package models

import "time"

type User struct {
	Id           uint64 `json:"id"`
	Fullname     string `json:"fullname"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Bio          string `json:"bio"`
}

type Session struct {
	UserId    uint64
	SessionId string
}

type Message struct {
	Id         uint64    `json:"id"`
	SenderId   uint64    `json:"senderId"`
	ReceiverId uint64    `json:"receiverId"`
	Text       string    `json:"text"`
	Timestamp  time.Time `json:"timestamp"`
}

type Conversation struct {
	UserId1   uint64    `json:"userId1"`
	UserId2   uint64    `json:"userId2"`
	Timestamp time.Time `json:"timestamp"`
}
