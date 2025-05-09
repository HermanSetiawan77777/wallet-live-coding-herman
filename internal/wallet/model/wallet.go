package model

type Wallet struct {
	ID      int `gorm:"primaryKey;identity(1,1);column:id"`
	UserID  int `gorm:"not null;column:user_id;index"`
	Balance int `gorm:"not null;default:0;column:balance"`
}

func (Wallet) TableName() string {
	return "wallets"
}