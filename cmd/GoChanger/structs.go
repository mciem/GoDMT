package main

import (
	"sync"

	"github.com/mciem/GoDMT/internal/utils"
)

type Config struct {
	Changer Changer `yaml:"changer"`
}

type Changer struct {
	Bios  []string `yaml:"bios"`
	Names []string `yaml:"names"`
}

var (
	wg sync.WaitGroup

	config  Config
	tokens  *utils.Cycle
	proxies *utils.Cycle
)
