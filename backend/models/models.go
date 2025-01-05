package models

import (
	"time"
)

type Role string

const (
	Admin  Role = "admin"
	Editor Role = "editor"
	Viewer Role = "viewer"
)

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Login      string `json:"login"`
	Role       Role   `json:"role"`
	TelegramID int64  `json:"telegram_id"`
}

type Quote struct {
	ID       string    `json:"id"`
	Text     string    `json:"text"`
	Author   string    `json:"author"`
	Datetime time.Time `json:"datetime"`
	Creator  User      `json:"creator"`
}
