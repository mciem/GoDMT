package main

import (
	"fmt"

	"github.com/mciem/GoDMT/internal/console"
	"github.com/mciem/GoDMT/internal/discord"
	"github.com/mciem/GoDMT/internal/utils"
)

func main() {
	token := "OTk3ODgxMTM3NjI2MTAzODg5.GUeJZl.ULeSQEA6oTKe01eY22KOgTDoGfCG9Mt4G90VBs"
	cID := "1142025546302230559"
	gID := "1142025545840873482"

	socket := discord.NewDiscordSocket(token, gID, cID)
	disc, err := discord.NewDiscord(token, "http://fwtvtpxaiftxiui57J-res-ROW-sid-33629637:gcgzywvkbdyo96I@gw.thunderproxies.net:5959")
	if err != nil {
		fmt.Println(err.Error())
	}

	socket.Run()

	users := socket.Users

	retries := 0
	for i, user := range users {
		if retries == 3 {
			return
		}

		s, _, err := disc.SendFriendRequest(user.ID, gID, cID)
		if utils.HandleError(err) {
			continue
		}

		if s {
			console.Log("SCC", "sent friend request to "+user.Username, map[string]string{
				"sent": fmt.Sprint(i - retries),
			})
		} else {
			console.Log("ERR", "failed to sent friend request to "+user.Username, map[string]string{
				"sent": fmt.Sprint(i - retries),
			})

			retries++
		}

	}
}
