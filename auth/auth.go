package auth

import (
	"encoding/json"
	"github.com/jiaying2001/agent/dto"
	"github.com/jiaying2001/agent/store"
	"github.com/jiaying2001/agent/transportation"
)

func Login(username, password string) bool {
	var str = dto.Credential{
		Username: username,
		Password: password,
	}
	body, _ := json.Marshal(&str)
	bytes := transportation.Post("/api/user/login", body)
	var response dto.Response
	json.Unmarshal(bytes, &response)
	if response.Code == 200 {
		store.Pass.AuthToken = response.Data
		store.Pass.UserName = str.Username
		return true
	} else {
		return false
	}
}
