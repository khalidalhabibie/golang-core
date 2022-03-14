package dingtalk

type Payload struct {
	At      AtModel   `json:"at"`
	Text    TextModel `json:"text"`
	Msgtype string    `json:"msgtype"`
}

type AtModel struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIDs []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

type TextModel struct {
	Content string `json:"content"`
}
