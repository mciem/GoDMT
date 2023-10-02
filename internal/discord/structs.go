package discord

import (
	"encoding/base64"
	"encoding/json"
)

type DisplayNamePayload struct {
	GlobalName string `json:"global_name"`
}

type BioPayload struct {
	Bio string `json:"bio"`
}

type SendMessagePayload struct {
	Content string `json:"content"`
	Tts     bool   `json:"tts"`
}

type CreateChannelPayload struct {
	Recipients []string `json:"recipients"`
}

type CreateChannelResponse struct {
	ID string `json:"id"`
}

type WebsocketSessionResponse struct {
	D D1 `json:"d"`
}

type D1 struct {
	Session_id string `json:"session_id"`
}

type WebsocketOnlinePayload struct {
	Op int `json:"op"`
	D  D   `json:"d"`
}

type Properties struct {
	Os      string `json:"$os"`
	Browser string `json:"$browser"`

	Device string `json:"$device"`
}
type Presence struct {
	Status     string `json:"status"`
	Since      int    `json:"since"`
	Activities []any  `json:"activities"`
	Afk        bool   `json:"afk"`
}

type D struct {
	Token      string     `json:"token"`
	Properties Properties `json:"properties"`
	Presence   Presence   `json:"presence"`
}

type InviteData struct {
	Type                     int      `json:"type"`
	Code                     string   `json:"code"`
	ExpiresAt                any      `json:"expires_at"`
	Guild                    GuildUno `json:"guild"`
	GuildID                  string   `json:"guild_id"`
	Channel                  Channel  `json:"channel"`
	ApproximateMemberCount   int      `json:"approximate_member_count"`
	ApproximatePresenceCount int      `json:"approximate_presence_count"`
}

type GuildUno struct {
	ID                       string   `json:"id"`
	Name                     string   `json:"name"`
	Splash                   any      `json:"splash"`
	Banner                   any      `json:"banner"`
	Description              any      `json:"description"`
	Icon                     string   `json:"icon"`
	Features                 []string `json:"features"`
	VerificationLevel        int      `json:"verification_level"`
	VanityURLCode            string   `json:"vanity_url_code"`
	NsfwLevel                int      `json:"nsfw_level"`
	Nsfw                     bool     `json:"nsfw"`
	PremiumSubscriptionCount int      `json:"premium_subscription_count"`
}

type Channel struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
}

type XContext struct {
	Location            string `json:"location"`
	LocationGuildID     string `json:"location_guild_id"`
	LocationChannelID   string `json:"location_channel_id"`
	LocationChannelType int    `json:"location_channel_type"`
}

type FriendRequestPayload struct {
	Username      string `json:"username"`
	Discriminator any    `json:"discriminator"`
}

func BuildXContext(invD InviteData) string {
	pd, _ := json.Marshal(XContext{
		Location:            "Join Guild",
		LocationGuildID:     invD.Guild.ID,
		LocationChannelID:   invD.Channel.ID,
		LocationChannelType: invD.Channel.Type,
	})

	return base64.RawStdEncoding.EncodeToString(pd)

}

var (
	headerOrder = []string{
		"accept",
		"accept-language",
		"content-type",
		"connection",
		"host",
		"origin",
		"referer",
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-platform",
		"sec-fetch-dest",
		"sec-fetch-mode",
		"sec-fetch-site",
		"user-agent",
		"x-context-properties",
		"x-debug-options",
		"x-discord-locale",
		"x-discord-timezone",
		"x-fingerprint",
		"x-super-properties",
		"x-track",
	}
)
