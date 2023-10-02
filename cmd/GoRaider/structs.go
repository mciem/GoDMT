package main

import (
	"sync"

	"github.com/mciem/GoDMT/internal/discord"
	"github.com/mciem/GoDMT/internal/utils"
)

type Config struct {
	Joiner Joiner `yaml:"joiner"`
	Raider Raider `yaml:"raider"`
}

type Joiner struct {
	Sleep  int    `yaml:"sleep"`
	Invite string `yaml:"invite"`
}

type Raider struct {
	Massping  bool   `yaml:"massping"`
	Message   string `yaml:"message"`
	ChannelID string `yaml:"channelID"`
	GuildID   string `yaml:"guildID"`
}

var (
	wg sync.WaitGroup

	config  Config
	tokens  *utils.Cycle
	proxies *utils.Cycle
	users   *utils.Cycle
	data    discord.InviteData
)
