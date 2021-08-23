package server

import (
	"fmt"
	"time"

	"github.com/dev-zipida-com/sos-ws/internal/chat"
	"github.com/gin-gonic/gin"
)

const HTML = "public/index.html"

type Ticker struct {
	UserCount int `json:"userCount"`
}

func Start() {
	h := chat.NewHub()
	go h.Run()

	router := gin.New()
	router.LoadHTMLFiles(HTML)

	router.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		chat.ServeWs(h, c.Writer, c.Request, roomId)
	})

	go func() {
		for {
			// log
			log := make([]string, 1)
			for c, connections := range h.GetRooms() {
				item := fmt.Sprintf("%s: %d ", c, len(connections))
				log = append(log, item)
			}

			fmt.Println(log)
			time.Sleep(time.Millisecond * 1000)
		}
	}()

	router.Run("0.0.0.0:18080")
}
