package model

import "time"

type CanvasHashLibrary struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Hash       string    `gorm:"column:hash;size:64;uniqueIndex" json:"hash"`
	SampleJSON string    `gorm:"column:sample_json;type:text" json:"sample_json"`
	SeenCount  uint64    `gorm:"column:seen_count;default:1" json:"seen_count"`
	FirstSeen  time.Time `gorm:"column:first_seen" json:"first_seen"`
	LastSeen   time.Time `gorm:"column:last_seen" json:"last_seen"`
}

func (CanvasHashLibrary) TableName() string { return "canvas_hash_library" }

type WebglHashLibrary struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Hash       string    `gorm:"column:hash;size:64;uniqueIndex" json:"hash"`
	SampleJSON string    `gorm:"column:sample_json;type:text" json:"sample_json"`
	SeenCount  uint64    `gorm:"column:seen_count;default:1" json:"seen_count"`
	FirstSeen  time.Time `gorm:"column:first_seen" json:"first_seen"`
	LastSeen   time.Time `gorm:"column:last_seen" json:"last_seen"`
}

func (WebglHashLibrary) TableName() string { return "webgl_hash_library" }

type FontsHashLibrary struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Hash       string    `gorm:"column:hash;size:64;uniqueIndex" json:"hash"`
	SampleJSON string    `gorm:"column:sample_json;type:text" json:"sample_json"`
	SeenCount  uint64    `gorm:"column:seen_count;default:1" json:"seen_count"`
	FirstSeen  time.Time `gorm:"column:first_seen" json:"first_seen"`
	LastSeen   time.Time `gorm:"column:last_seen" json:"last_seen"`
}

func (FontsHashLibrary) TableName() string { return "fonts_hash_library" }
