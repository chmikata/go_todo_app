package entity

import "time"

type UserId int64

type User struct {
	ID       UserId    `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Password string    `json:"password" db:"password"`
	Role     string    `json:"role" db:"role"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modifie"`
}
