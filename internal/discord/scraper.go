package discord

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mciem/GoDMT/internal/utils"
)

type User struct {
	Avatar               string `json:"avatar"`
	AvatarDecorationData any    `json:"avatar_decoration_data"`
	Bot                  bool   `json:"bot"`
	Discriminator        string `json:"discriminator"`
	DisplayName          string `json:"display_name"`
	GlobalName           string `json:"global_name"`
	ID                   string `json:"id"`
	PublicFlags          int    `json:"public_flags"`
	Username             string `json:"username"`
}

type UserProfile struct {
	Bio string `json:"bio"`
}

type DiscordSocket struct {
	MAX_ITER         int
	Token            string
	GuildID          string
	ChannelID        string
	BlacklistedRoles []string
	BlacklistedUsers []string
	SocketHeaders    map[string]string
	EndScraping      bool
	Statuses         []string
	Guilds           map[string]map[string]int
	Users            []User
	Ranges           [][]int
	LastRange        int
	PacketsReceived  int
	Msgs             []int
	D                int
	Iter             int
	BigIter          int
	Finished         bool
	WS               *websocket.Conn
}

type MemberData struct {
	OnlineCount  int           `json:"online_count"`
	MemberCount  int           `json:"member_count"`
	ID           string        `json:"id"`
	GuildID      string        `json:"guild_id"`
	HoistedRoles []string      `json:"hoisted_roles"`
	Types        []string      `json:"types"`
	Locations    []int         `json:"locations"`
	Updates      []interface{} `json:"updates"`
}

type WebSocketMessage struct {
	Op int         `json:"op"`
	D  interface{} `json:"d"`
	T  string      `json:"t"`
}

func NewDiscordSocket(token, guildID, channelID string) *DiscordSocket {
	dialer := websocket.Dialer{}
	ws, _, err := dialer.Dial("wss://gateway.discord.gg/?encoding=json&v=9", http.Header{
		"Accept-Language":          []string{"en-US,en;q=0.9"},
		"Cache-Control":            []string{"no-cache"},
		"Pragma":                   []string{"no-cache"},
		"Sec-WebSocket-Extensions": []string{"permessage-deflate; client_max_window_bits"},
		"User-Agent":               []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
	})

	if utils.HandleError(err) {
		return &DiscordSocket{}
	}

	return &DiscordSocket{
		MAX_ITER:        10,
		Token:           token,
		GuildID:         guildID,
		ChannelID:       channelID,
		Guilds:          make(map[string]map[string]int),
		Statuses:        []string{},
		Ranges:          [][]int{{0}},
		LastRange:       0,
		PacketsReceived: 0,
		Msgs:            make([]int, 0),
		D:               1,
		Iter:            0,
		BigIter:         0,
		Finished:        false,
		WS:              ws,
	}
}

func (ds *DiscordSocket) getRanges(index int, multiplier float64, memberCount int) [][]int {
	initialNum := int(float64(index) * multiplier)
	rangesList := [][]int{{initialNum, initialNum + 99}}
	if memberCount > initialNum+99 {
		rangesList = append(rangesList, []int{initialNum + 100, initialNum + 199})
	}
	if !containsRange(rangesList, []int{0, 99}) {
		rangesList = append([][]int{{0, 99}}, rangesList...)
	}
	return rangesList
}

func containsRange(ranges [][]int, targetRange []int) bool {
	for _, r := range ranges {
		if r[0] == targetRange[0] && r[1] == targetRange[1] {
			return true
		}
	}
	return false
}

func (ds *DiscordSocket) parseGuildMemberListUpdate(response *WebSocketMessage) *MemberData {
	data := response.D.(map[string]interface{})
	memberdata := &MemberData{
		OnlineCount: int(data["online_count"].(float64)),
		MemberCount: int(data["member_count"].(float64)),
		ID:          data["id"].(string),
		GuildID:     data["guild_id"].(string),
		//HoistedRoles: data["groups"].([]string),
		Types:     []string{},
		Locations: []int{},
		Updates:   []interface{}{},
	}
	ops := data["ops"].([]interface{})
	for _, chunk := range ops {
		op := chunk.(map[string]interface{})
		memberdata.Types = append(memberdata.Types, op["op"].(string))
		if op["op"].(string) == "SYNC" || op["op"].(string) == "INVALIDATE" {
			f := op["range"].([]interface{})

			for _, x := range f {
				memberdata.Locations = append(memberdata.Locations, int(x.(float64)))
			}
			if op["op"].(string) == "SYNC" {
				memberdata.Updates = append(memberdata.Updates, op["items"])
			} else {
				memberdata.Updates = append(memberdata.Updates, []interface{}{})
			}
		} else if op["op"].(string) == "INSERT" || op["op"].(string) == "UPDATE" || op["op"].(string) == "DELETE" {
			memberdata.Locations = append(memberdata.Locations, int(op["index"].(float64)))
			if op["op"].(string) == "DELETE" {
				memberdata.Updates = append(memberdata.Updates, []interface{}{})
			} else {
				memberdata.Updates = append(memberdata.Updates, op["item"])
			}
		}
	}
	return memberdata
}

func (ds *DiscordSocket) findMostReoccuring(list []int) int {
	counts := make(map[int]int)
	for _, value := range list {
		counts[value]++
	}
	maxCount := 0
	maxValue := 0
	for key, count := range counts {
		if count > maxCount {
			maxCount = count
			maxValue = key
		}
	}
	return maxValue
}

func (ds *DiscordSocket) Run() []User {
	ds.sockOpen()

	for !ds.Finished {
		_, m, err := ds.WS.ReadMessage()

		if err != nil {
			return []User{}
		}

		ds.sockMessage(ds.WS, string(m))
	}

	return ds.Users
}

func (ds *DiscordSocket) scrapeUsers() {
	if !ds.EndScraping {
		ds.send(fmt.Sprintf(`{"op":14,"d":{"guild_id":"%s","typing":true,"activities":true,"threads":true,"channels":{"%s":%s}}}`, ds.GuildID, ds.ChannelID, toJSON(ds.Ranges)))
	}
}

func (ds *DiscordSocket) sockOpen() {
	ds.send(`{"op":2,"d":{"token":"` + ds.Token + `","capabilities":125,"properties":{"os":"Windows","browser":"Firefox","device":"","system_locale":"it-IT","browser_user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:94.0) Gecko/20100101 Firefox/94.0","browser_version":"94.0","os_version":"10","referrer":"","referring_domain":"","referrer_current":"","referring_domain_current":"","release_channel":"stable","client_build_number":232643,"client_event_source":null},"presence":{"status":"online","since":0,"activities":[],"afk":false},"compress":false,"client_state":{"guild_hashes":{},"highest_last_message_id":"0","read_state_version":0,"user_guild_settings_version":-1,"user_settings_version":-1}}}`)
}

func (ds *DiscordSocket) heartbeatThread(interval time.Duration) {
	for {
		ds.send(fmt.Sprintf(`{"op":1,"d":%d}`, ds.PacketsReceived))
		time.Sleep(interval)
	}
}

func (ds *DiscordSocket) sockMessage(ws *websocket.Conn, message string) {
	var decoded WebSocketMessage
	err := json.Unmarshal([]byte(message), &decoded)
	if err != nil {
		return
	}

	if decoded.Op != 11 {
		ds.PacketsReceived++
	}

	if decoded.Op == 10 {
		go ds.heartbeatThread(time.Duration(41250/1000) * time.Second)
	}

	if decoded.T == "READY" {
		for _, guild := range decoded.D.(map[string]interface{})["guilds"].([]interface{}) {
			guildData := guild.(map[string]interface{})
			ds.Guilds[guildData["id"].(string)] = map[string]int{"member_count": int(guildData["member_count"].(float64))}
		}
	}

	if decoded.T == "READY_SUPPLEMENTAL" {
		ds.Ranges = ds.getRanges(0, 100, ds.Guilds[ds.GuildID]["member_count"])
		ds.scrapeUsers()

	} else if decoded.T == "GUILD_MEMBER_LIST_UPDATE" {
		parsed := ds.parseGuildMemberListUpdate(&decoded)
		ds.Msgs = append(ds.Msgs, len(ds.Users))

		if ds.D == len(ds.Users) {
			ds.Iter++
			if ds.Iter == ds.MAX_ITER {
				ds.Finished = true
				return
			}
		}
		ds.D = ds.findMostReoccuring(ds.Msgs)

		if parsed.GuildID == ds.GuildID && (containsString(parsed.Types, "SYNC") || containsString(parsed.Types, "UPDATE")) {
			for elem, index := range parsed.Types {
				if index == "SYNC" {
					for _, item := range parsed.Updates[elem].([]interface{}) {
						if item.(map[string]interface{})["member"] == nil {
							continue
						}

						member := item.(map[string]interface{})["member"].(map[string]interface{})
						//obj := map[string]interface{}{
						//	"tag": member["user"].(map[string]interface{})["username"].(string) + "#" + member["user"].(map[string]interface{})["discriminator"].(string),
						//	"id":  member["user"].(map[string]interface{})["id"].(string),
						//}

						activites := member["presence"].(map[string]interface{})["activities"].([]interface{})
						if len(activites) != 0 {
							for _, x := range activites {
								if x.(map[string]interface{})["id"] == "custom name" {
									ds.Statuses = append(ds.Statuses, x.(map[string]interface{})["state"].(string))
								}
							}

						}

						if member["user"].(map[string]interface{})["bot"] == false {
							r, _ := json.Marshal(member["user"])

							var user User
							json.Unmarshal(r, &user)

							ds.Users = append(ds.Users, user)
						}
					}
				} else if index == "UPDATE" {
					for _, item := range parsed.Updates[elem].(map[string]interface{}) {
						if item.(map[string]interface{})["member"] == nil {
							continue
						}

						member := item.(map[string]interface{})["member"].(map[string]interface{})
						//obj := map[string]interface{}{
						//	"tag": member["user"].(map[string]interface{})["username"].(string) + "#" + member["user"].(map[string]interface{})["discriminator"].(string),
						//	"id":  member["user"].(map[string]interface{})["id"].(string),
						//}

						activites := member["presence"].(map[string]interface{})["activities"].([]interface{})
						if len(activites) != 0 {
							for _, x := range activites {
								if x.(map[string]interface{})["id"] == "custom name" {
									ds.Statuses = append(ds.Statuses, x.(map[string]interface{})["state"].(string))
								}
							}

						}

						if member["user"].(map[string]interface{})["bot"] == false {
							r, _ := json.Marshal(member["user"])

							var user User
							json.Unmarshal(r, &user)

							ds.Users = append(ds.Users, user)
						}
					}
				}
				ds.LastRange++
				ds.Ranges = ds.getRanges(ds.LastRange, 100, ds.Guilds[ds.GuildID]["member_count"])
				time.Sleep(450 * time.Millisecond)
				ds.scrapeUsers()
			}
		}
	}
}

func (ds *DiscordSocket) send(message string) {
	err := ds.WS.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return
	}
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func containsString(list []string, target string) bool {
	for _, str := range list {
		if str == target {
			return true
		}
	}
	return false
}
