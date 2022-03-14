package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func (d *DingTalk) SendtoDingtalk(accessToken, secret, content string, atMobiles, atUserIds []string, isAtAll bool) error {
	timeStamp, sign := getSignedKey(secret)

	urlDingtalk := os.Getenv("DINGTALK_URL")
	headersURL := os.Getenv("DINGTALK_HEADERS")

	dataPost := fmt.Sprintf(
		"%s/robot/send?access_token=%s&timestamp=%s&sign=%s",
		urlDingtalk, accessToken, timeStamp, sign)

	payload := Payload{}

	// at
	payload.At.AtMobiles = atMobiles
	payload.At.AtUserIDs = atUserIds
	payload.At.IsAtAll = isAtAll

	// text
	payload.Text.Content = content

	// msgtype
	payload.Msgtype = "text"

	postBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("An Error marshal post body, err ", err)
		return err
	}

	responseBody := bytes.NewBuffer(postBody)

	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(dataPost, headersURL, responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occurred %v", err)
		return err
	}
	defer resp.Body.Close()

	return nil

}

func getSignedKey(secret string) (string, string) {
	timeNow := fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond))
	secretEnc := []byte(secret)

	stringToSign := fmt.Sprintf("%v\n%v", timeNow, secret)

	stringToSignEnc := []byte(stringToSign)
	hmacGenerate := hmac.New(sha256.New, secretEnc)
	hmacGenerate.Write(stringToSignEnc)

	sign := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(
		hmacGenerate.Sum(nil),
	)))

	return timeNow, sign

}
