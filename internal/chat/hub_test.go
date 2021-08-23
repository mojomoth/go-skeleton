package chat

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
)

func TestFind(t *testing.T) {
	r := regexp.MustCompile("\\(\\/.*?\\)")
	result := r.FindString("(/AAA) aasdasdaxcz (/zxcssdasd)")
	fmt.Println(result)

	result = r.FindString("aasdasdaxcz")
	fmt.Println(result)

	if result != "" {
		fmt.Println("Found")
	} else {
		fmt.Println("Not Found")
	}
}

func TestChatbot(t *testing.T) {
	bubbleData := BubbleData{
		Description: "(/도와줘SOS)",
	}

	bubble := Bubble{
		Type: "text",
		Data: bubbleData,
	}

	requestBody := RequestBody{
		Version:   "v2",
		UserId:    "test",
		Timestamp: 123456,
		Bubbles: []Bubble{
			bubble,
		},
		Event: "send",
	}

	secretKey := "WHd2RkZIWmFNVGFVcFpyZVhSekVZV0ppd2xxUE9zSUE="
	url := "https://6e875d362ca2402d8936c323c582701e.apigw.ntruss.com/custom/v1/5268/f6ccaa7ac017166dddc86d1b4d4f102653f7508a12365aaf50b913aee51b4a56"

	h := hmac.New(sha256.New, []byte(secretKey))
	b, _ := json.Marshal(requestBody)
	h.Write(b)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(signature)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-NCP-CHATBOT_SIGNATURE", signature)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Body:", string(body))
}
