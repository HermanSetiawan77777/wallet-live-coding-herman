package model

import "time"

type Transaction struct {
    ID        int       `gorm:"primaryKey;identity(1,1);column:id"`
    UserID    int       `gorm:"not null;column:user_id;index"`
    Type      string    `gorm:"not null;column:type;size:50"`
    Amount    int       `gorm:"not null;column:amount"`
    CreatedAt time.Time `gorm:"not null;column:created_at;default:getdate()"`
}

func (Transaction) TableName() string {
    return "transactions"
}