package config

import (
	"encoding/json"
	"fmt"
	"github.com/jiaying2001/agent/harvester"
	"github.com/jiaying2001/agent/transportation"
)

type HarvesterConfigResp struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Configs []harvester.Harvester `json:"data"`
}

func GetConfigs() *[]harvester.Harvester {
	configsJson := transportation.Get("/api/harvester", nil)
	var configs HarvesterConfigResp
	// 将JSON数据解码到people数组中
	if err := json.Unmarshal(configsJson, &configs); err != nil {
		fmt.Println("解码JSON时出错:", err)
		return nil
	}
	if configs.Code != 200 {
		fmt.Println("请求状态码异常")
	}
	return &configs.Configs
}
