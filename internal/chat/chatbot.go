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
	"time"
)

type BubbleData struct {
	Description string `json:"description"`
}

type Bubble struct {
	Type string     `json:"type"`
	Data BubbleData `json:"data"`
}

type RequestBody struct {
	Version   string   `json:"version"`
	UserId    string   `json:"userId"`
	Timestamp int64    `json:"timestamp"`
	Bubbles   []Bubble `json:"bubbles"`
	Event     string   `json:"event"`
}

type ContentTableDataDataActionData struct {
	Postback     string `json:"postback"`
	PostbackFull string `json:"postbackFull"`
}

type ContentTableDataDataAction struct {
	Type string                         `json:"type"`
	Data ContentTableDataDataActionData `json:"data"`
}

type ContentTableDataData struct {
	Type   string                     `json:"type"`
	Action ContentTableDataDataAction `json:"action"`
}

type ContentTableData struct {
	Type  string               `json:"type"`
	Title string               `json:"title"`
	Data  ContentTableDataData `json:"data"`
}

type ContentTable struct {
	RowSpan int              `json:"rowSpan"`
	ColSpan int              `json:"colSpan"`
	Data    ContentTableData `json:"data"`
}

type ResponseBubbleData struct {
	Cover        Bubble           `json:"cover"`
	ContentTable [][]ContentTable `json:"contentTable"`
}

type ResponseBubble struct {
	Type string             `json:"type"`
	Data ResponseBubbleData `json:"data"`
}

type ResponseBody struct {
	Version   string           `json:"version"`
	UserId    string           `json:"userId"`
	Timestamp int64            `json:"timestamp"`
	Bubbles   []ResponseBubble `json:"bubbles"`
	Event     string           `json:"event"`
}

type RequestDatabaseBody struct {
	Name        string `json:"name"`
	Message     string `json:"message"`
	Client      string `json:"client"`
	IsRead      bool   `json:"isRead"`
	MessageType string `json:"messageType"`
	ChatRoom    string `json:"chatRoom"`
	CreatedAt   string `json:"createdAt"`
}

const (
	SECRET_KEY   = "WHd2RkZIWmFNVGFVcFpyZVhSekVZV0ppd2xxUE9zSUE="
	URL          = "https://6e875d362ca2402d8936c323c582701e.apigw.ntruss.com/custom/v1/5268/f6ccaa7ac017166dddc86d1b4d4f102653f7508a12365aaf50b913aee51b4a56"
	CHATBOT_ID   = "chatbot"
	CHATBOT_NAME = "SOS봇"
)

const (
	API_URL = "http://3.36.180.31:11337/chats"
)

const (
	INPUT_FORMAT  = "2006-01-02T15:04:05.999999999-07:00"
	OUTPUT_FORMAT = "2006-01-02T15:04:05.000Z"
)

func timestampToJavaScriptISO(s string) (string, error) {
	t, err := time.Parse(INPUT_FORMAT, s)
	if err != nil {
		return "", err
	}
	return t.UTC().Format(OUTPUT_FORMAT), nil
}

func chatbot(send chan<- []byte, command string, chatroom string) {
	timestamp := time.Now().UnixNano() / 1000000
	bubbleData := BubbleData{
		Description: command,
	}

	bubble := Bubble{
		Type: "text",
		Data: bubbleData,
	}

	requestBody := RequestBody{
		Version:   "v2",
		UserId:    CHATBOT_ID,
		Timestamp: timestamp,
		Bubbles: []Bubble{
			bubble,
		},
		Event: "send",
	}

	h := hmac.New(sha256.New, []byte(SECRET_KEY))
	b, _ := json.Marshal(requestBody)
	h.Write(b)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(b))
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

	res := ResponseBody{}
	json.Unmarshal(body, &res)

	createdAt := time.Now().Format(time.RFC3339)
	c := Chat{
		Name:        CHATBOT_NAME,
		Message:     string(body),
		Client:      CHATBOT_ID,
		CreatedAt:   createdAt,
		IsRead:      false,
		MessageType: "chatbot",
		ChatRoom:    chatroom,
	}

	chatBytes, _ := json.Marshal(c)
	send <- chatBytes

	saveChat(string(body), chatroom, createdAt)
}

func saveChat(message string, chatroom string, createdAt string) {
	requestBody := RequestDatabaseBody{
		Name:        CHATBOT_NAME,
		Message:     message,
		IsRead:      false,
		MessageType: "chatbot",
		ChatRoom:    chatroom,
		CreatedAt:   createdAt,
	}

	b, _ := json.Marshal(requestBody)
	fmt.Println(string(b))
	req, _ := http.NewRequest("POST", API_URL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("strapi response Status:", resp.Status)
	fmt.Println("strapi response Body:", string(body))
}
