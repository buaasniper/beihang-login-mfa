package model

type FingerprintLog struct {
	ID          uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string   `json:"username"`
	Fingerprint string   `gorm:"type:text" json:"fingerprint"`
	URL         string   `gorm:"column:url" json:"url"`
	DeltaTime   *float64 `gorm:"column:delta_time" json:"delta_time"`
	ClickTime   string   `gorm:"column:click_time" json:"click_time"`
}

func (FingerprintLog) TableName() string { return "fingerprint_logs" }
