package main

import (
	"github.com/mciem/GoDMT/internal/console"
	"github.com/mciem/GoDMT/internal/discord"
	"github.com/mciem/GoDMT/internal/utils"
)

var (
	proxie *utils.Cycle

	ersdsdsr error
)

func main() {
	proxie, ersdsdsr = utils.NewFromFile("../../assets/proxies.txt")
	if utils.HandleError(ersdsdsr) {
		return
	}

	pr, _ := proxie.Next()
	dc, err := discord.NewDiscord("MTE1MjY5NjE2NDcxODc1MTg0Ng.GKvtm4.8XePlCY6c7JmfKb3ZoKaMvO2xppgzLhSWqRIsc", "http://"+pr)
	if utils.HandleError(err) {
		return
	}

	s, x, ers := dc.JoinServer("REdf7sVH")
	if utils.HandleError(ers) {
		return
	}

	if s {
		console.Log("SCC", x, map[string]string{
			"token": dc.Token,
		})
	} else {
		console.Log("FLD", x, map[string]string{
			"token":  dc.Token,
			"reason": x,
		})
	}

}
