package model

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"column:username;type:varchar(128);uniqueIndex;not null" json:"username"`
	Password string `gorm:"column:password;type:varchar(255);not null" json:"password"`
}

func (User) TableName() string { return "users" }
