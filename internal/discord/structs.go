package discord

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
	Type                     int     `json:"type"`
	Code                     string  `json:"code"`
	ExpiresAt                any     `json:"expires_at"`
	Guild                    Guild   `json:"guild"`
	GuildID                  string  `json:"guild_id"`
	Channel                  Channel `json:"channel"`
	ApproximateMemberCount   int     `json:"approximate_member_count"`
	ApproximatePresenceCount int     `json:"approximate_presence_count"`
}

type Guild struct {
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
