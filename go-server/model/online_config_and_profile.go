package model

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

var validRiskHashTypes = map[string]struct{}{
	"cookie": {},
	"canvas": {},
	"webgl":  {},
	"fonts":  {},
}

var sha256HexRe = regexp.MustCompile("^[a-f0-9]{64}$")

type OnlineRuleConfig struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RuleKey    string    `gorm:"column:rule_key;size:64;uniqueIndex" json:"rule_key"`
	Enabled    bool      `gorm:"column:enabled;default:true;index" json:"enabled"`
	Priority   int       `gorm:"column:priority;default:100;index" json:"priority"`
	Source     string    `gorm:"column:source;size:16" json:"source"`
	FieldPath  string    `gorm:"column:field_path;size:128" json:"field_path"`
	Operator   string    `gorm:"column:operator;size:32" json:"operator"`
	Value      string    `gorm:"column:value;type:text" json:"value"`
	Decision   string    `gorm:"column:decision;size:16;index" json:"decision"`
	Version    string    `gorm:"column:version;size:32;index" json:"version"`
	UpdatedAt  time.Time `gorm:"column:updated_at;index" json:"updated_at"`
}

func (OnlineRuleConfig) TableName() string { return "online_rule_config" }

type RiskAccountProfile struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"column:username;size:128;uniqueIndex" json:"username"`
	RiskLevel string    `gorm:"column:risk_level;size:16;index" json:"risk_level"`
	RiskScore int       `gorm:"column:risk_score;default:0" json:"risk_score"`
	Reason    string    `gorm:"column:reason;type:text" json:"reason"`
	UpdatedAt time.Time `gorm:"column:updated_at;index" json:"updated_at"`
}

func (RiskAccountProfile) TableName() string { return "risk_account_profile" }

type RiskUAProfile struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UAHash      string    `gorm:"column:ua_hash;size:64;uniqueIndex" json:"ua_hash"`
	UA          string    `gorm:"column:ua;type:text" json:"ua"`
	RiskLevel   string    `gorm:"column:risk_level;size:16;index" json:"risk_level"`
	RiskScore   int       `gorm:"column:risk_score;default:0" json:"risk_score"`
	SignalCount uint64    `gorm:"column:signal_count;default:0" json:"signal_count"`
	UpdatedAt   time.Time `gorm:"column:updated_at;index" json:"updated_at"`
}

func (RiskUAProfile) TableName() string { return "risk_ua_profile" }

type RiskHashProfile struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	HashType  string    `gorm:"column:hash_type;size:16;uniqueIndex:uq_hash_type_value,priority:1;index" json:"hash_type"`
	HashValue string    `gorm:"column:hash_value;size:64;uniqueIndex:uq_hash_type_value,priority:2" json:"hash_value"`
	RiskLevel string    `gorm:"column:risk_level;size:16;index" json:"risk_level"`
	RiskScore int       `gorm:"column:risk_score;default:0" json:"risk_score"`
	Reason    string    `gorm:"column:reason;type:text" json:"reason"`
	UpdatedAt time.Time `gorm:"column:updated_at;index" json:"updated_at"`
}

func (RiskHashProfile) TableName() string { return "risk_hash_profile" }

func (r *RiskHashProfile) BeforeCreate(tx *gorm.DB) error {
	return r.validate()
}

func (r *RiskHashProfile) BeforeUpdate(tx *gorm.DB) error {
	return r.validate()
}

func (r *RiskHashProfile) validate() error {
	r.HashType = strings.ToLower(strings.TrimSpace(r.HashType))
	r.HashValue = strings.ToLower(strings.TrimSpace(r.HashValue))

	if _, ok := validRiskHashTypes[r.HashType]; !ok {
		return errors.New("invalid hash_type, allowed values: cookie, canvas, webgl, fonts")
	}
	if !sha256HexRe.MatchString(r.HashValue) {
		return errors.New("invalid hash_value, expected 64-char lowercase sha256 hex")
	}
	return nil
}
