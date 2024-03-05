package dto

type Message struct {
	Header   Header   `json:"header"`
	Body     Body     `json:"body"`
	Property Property `json:"property"`
}

type Header struct {
	TraceId   string `json:"traceId"`
	AuthToken string `json:"authToken"`
	Os        string `json:"os"`
	AppName   string `json:"appName"`
	Pid       string `json:"pid"`
	UserID    int    `json:"userID"`
}

type Body struct {
	Content string `json:"content"`
}

type Property struct {
	StartTimeStamp []int64 `json:"startTimestamp"`
	EndTimeStamp   []int64 `json:"endTimestamp"`
}
