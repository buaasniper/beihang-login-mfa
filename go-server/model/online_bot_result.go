package model

import "time"

type OnlineBotResult struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BfpEventID  uint      `gorm:"column:bfp_event_id;index" json:"bfp_event_id"`
	Username    string    `gorm:"column:username;size:128;index" json:"username"`
	Decision    string    `gorm:"column:decision;size:16;index" json:"decision"`
	ReasonsJSON string    `gorm:"column:reasons_json;type:text" json:"reasons_json"`
	CreatedAt   time.Time `gorm:"column:created_at;index" json:"created_at"`
}

func (OnlineBotResult) TableName() string { return "online_bot_result" }
