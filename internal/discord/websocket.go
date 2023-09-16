package discord

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

func (d *Discord) OnlineToken() (string, error) {
	dialer := websocket.Dialer{}
	ws, _, err := dialer.Dial("wss://gateway.discord.gg/?encoding=json&v=9", http.Header{
		"Origin":     []string{"https://discord.com"},
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
	})

	if err != nil {
		return "", err
	}

	pd, er := json.Marshal(WebsocketOnlinePayload{
		Op: 2,
		D: D{
			Token: d.Token,
			Properties: Properties{
				Os:      "windows",
				Browser: "chrome",
				Device:  "pc",
			},
			Presence: Presence{
				Status:     "online",
				Since:      0,
				Activities: []any{},
				Afk:        false,
			},
		},
	})

	if er != nil {
		return "", er
	}

	ers := ws.WriteMessage(websocket.TextMessage, pd)
	if ers != nil {
		return "", ers
	}

	for i := 0; i < 10; i++ {
		_, rec, err := ws.ReadMessage()
		if err != nil {
			return "", err
		}

		var resp WebsocketSessionResponse
		json.Unmarshal(rec, &resp)

		if resp.D.Session_id != "" {
			return resp.D.Session_id, nil
		}
	}

	return "", nil

}
