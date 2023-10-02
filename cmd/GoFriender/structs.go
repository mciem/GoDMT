package main

import (
	"sync"

	"github.com/mciem/GoDMT/internal/discord"
	"github.com/mciem/GoDMT/internal/utils"
)

type Config struct {
	Friender Friender `yaml:"friender"`
	Joiner   Joiner   `yaml:"joiner"`
}

type Friender struct {
	Retries     int      `yaml:"retries"`
	Sleep       int      `yaml:"sleep"`
	Invite      string   `yaml:"invite"`
	BlacklistID []string `yaml:"blacklistID"`
}

type Joiner struct {
	Sleep int `yaml:"sleep"`
}

var (
	wg sync.WaitGroup

	config  Config
	tokens  *utils.Cycle
	proxies *utils.Cycle
	data    discord.InviteData
	usrs    []discord.User
)
