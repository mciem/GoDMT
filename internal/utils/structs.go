package utils

type XContext struct {
	Location            string `json:"location"`
	LocationGuildID     string `json:"location_guild_id"`
	LocationChannelID   string `json:"location_channel_id"`
	LocationChannelType int    `json:"location_channel_type"`
}
