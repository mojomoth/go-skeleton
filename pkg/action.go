package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/dev-zipida-com/sos-detection-to-protection/internal"
)

const (
	CLIENTS                 = "/clients"
	DETECTIONS_COUNT        = "/detections/count?client="
	PROTECTIONS             = "/protections"
	PROTECTIONS_FROM_CLIENT = "/protections?client="
)

const (
	DEFAULT_WEB_COUNT          = 500
	DEFAULT_SEARCHED_WEB_COUNT = 500
)

type Client struct {
	Id  string `json:"_id"`
	Uid string `json:"uid"`
}

type Detections struct {
	WebUrl        string `json:"web_url"`
	WebName       string `json:"web_name"`
	WebChannel    string `json:"web_channel"`
	WebCategory   string `json:"web_category"`
	PageTitle     string `json:"page_title"`
	PageSource    string `json:"page_source"`
	PageUrl       string `json:"page_url"`
	PageCategory  string `json:"page_category"`
	SnapshotUrl   string `json:"snapshot_url"`
	OriginFaceUrl string `json:"origin_face_url"`
	FoundFaceUrl  string `json:"found_face_url"`
	VideoUrl      string `json:"video_url"`
	ImageUrl      string `json:"image_url"`
	Client        Client `json:"client"`
}

type Protection struct {
	Id               string `json:"_id"`
	ProtectedCount   int    `json:"protectedCount"`
	WebCount         int    `json:"webCount"`
	SearchedWebCount int    `json:"searchedWebCount"`
	Client           string `json:"client"`
}

func Start() {
	// get clients
	res, _, err := internal.Call("GET", CLIENTS, nil)
	if err != nil {
		log.Fatalln(err)
	}

	clients := make([]Client, 0)
	json.Unmarshal(res, &clients)

	// loop
	for _, client := range clients {
		// get count
		countBody, _, err := internal.Call("GET", DETECTIONS_COUNT+client.Id, nil)
		if err != nil {
			log.Fatalln(err)
		}

		// check count
		count, _ := strconv.Atoi(string(countBody))
		if count == 0 {
			continue
		}

		protectionsBody, _, err := internal.Call("GET", PROTECTIONS_FROM_CLIENT+client.Id, nil)
		if err != nil {
			log.Fatalln(err)
		}

		protections := make([]Protection, 0)
		json.Unmarshal(protectionsBody, &protections)

		protection := Protection{
			ProtectedCount:   count,
			WebCount:         DEFAULT_WEB_COUNT,
			SearchedWebCount: DEFAULT_SEARCHED_WEB_COUNT,
			Client:           client.Id,
		}

		protectionBytes, _ := json.Marshal(protection)

		// update and insert
		if len(protections) == 0 {
			fmt.Println("POST", client.Id)
			internal.Call("POST) client id : ", PROTECTIONS, protectionBytes)
		} else {
			fmt.Println("PUT", protections[0].Id)
			internal.Call("PUT) protection id : ", PROTECTIONS+"/"+protections[0].Id, protectionBytes)
		}
	}
}
