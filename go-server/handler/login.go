package handler

import (
	"encoding/json"
	"net/http"

	"go-server/model"

	"gorm.io/gorm"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, map[string]any{"success": false, "message": "请求格式错误"})
			return
		}

		var user model.User
		result := db.Where("username = ?", req.Username).First(&user)

		if result.Error != nil {
			writeJSON(w, map[string]any{"success": false, "message": "用户不存在"})
			return
		}

		if user.Password != req.Password {
			writeJSON(w, map[string]any{"success": false, "message": "密码错误"})
			return
		}

		writeJSON(w, map[string]any{"success": true, "message": "登录成功"})
	}
}

func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, map[string]any{"success": false, "message": "请求格式错误"})
			return
		}

		var existing model.User
		if db.Where("username = ?", req.Username).First(&existing).Error == nil {
			writeJSON(w, map[string]any{"success": false, "message": "用户名已存在"})
			return
		}

		newUser := model.User{
			Username: req.Username,
			Password: req.Password,
		}
		if err := db.Create(&newUser).Error; err != nil {
			writeJSON(w, map[string]any{"success": false, "message": "注册失败"})
			return
		}

		writeJSON(w, map[string]any{"success": true, "message": "注册成功，请去登录"})
	}
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}
