package platform

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

func getSignedKey(secret string) (string, string) {
	timeNow := time.Now().String()
	secretEnc := []byte(secret)

	stringToSign := fmt.Sprintf("%v\n%v", timeNow, secret)

	stringToSignEnc := []byte(stringToSign)

	hmacGenerate := hmac.New(sha256.New, secretEnc)
	hmacCode, err := hmacGenerate.Write(stringToSignEnc)
	if err != nil {
		log.Println("error parsing hmac: ", err)

	}

	sign := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(hmacCode))))

	return timeNow, sign

}

func dingtalk(accessToken, secret, content string, atMobiles, atUserIds []string, isAtAll bool) {

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
		return
	}

	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(dataPost, headersURL, responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occurred %v", err)
		return
	}
	defer resp.Body.Close()

}
