package model

type User struct {
	ID       int    `gorm:"primaryKey;identity(1,1);column:id"`
	Username string `gorm:"not null;column:username;index"`
	Email    string `gorm:"not null;column:email;index"`
}

func (User) TableName() string {
	return "users"
}