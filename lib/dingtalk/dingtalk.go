package dingtalk

// var env
// - DINGTALK_HEADERS
// - DINGTALK_URL

type DingTalk struct {
}

type Client interface {
	SendtoDingtalk(accessToken, secret, content string, atMobiles, atUserIds []string, isAtAll bool) error
}

func NewClient() Client {
	return &DingTalk{}
}
