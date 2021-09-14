package models

import "time"

// User Модель пользователя
type User struct {
	ID       int64
	Login    string
	Password string
	Name     string
	Birth    time.Time
	Balance  int64
	Created  time.Time
}
