package dto

import (
	"encoding/json"
	"time"
)

type LogMessage struct {
	Message   string    `json:"message"`
	Path      string    `json:"path"`
	AuthToken string    `json:"authToken"`
	Timestamp time.Time `json:"@timestamp"`
}

func BuildLogMessage(msg, path, authToken string, timestamp time.Time) []byte {
	bytes, _ := json.Marshal(&LogMessage{
		msg,
		path,
		authToken,
		timestamp,
	})
	return bytes
}
