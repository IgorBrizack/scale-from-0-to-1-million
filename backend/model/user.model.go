package model

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	SecondName string    `json:"second_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users_data"
}
