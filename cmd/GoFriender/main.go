package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/mciem/GoDMT/internal/console"
	"github.com/mciem/GoDMT/internal/discord"
	"github.com/mciem/GoDMT/internal/utils"
	"gopkg.in/yaml.v2"
)

func checkIfInBlacklist(x string) bool {
	for _, id := range config.Friender.BlacklistID {
		if id == x {
			return true
		}
	}

	return false
}

func main() {
	file, err := ioutil.ReadFile("config.yaml")
	if utils.HandleError(err) {
		return
	}

	yaml.Unmarshal(file, &config)

	tokens, _ = utils.NewFromFile("../../assets/tokens.txt")
	proxies, _ = utils.NewFromFile("../../assets/proxies.txt")

	xx := 0
	for i, token := range tokens.List {
		sent := 0

		prox, _ := proxies.Next()
		disc, err := discord.NewDiscord(token, "http://"+prox)
		if utils.HandleError(err) {
			continue
		}

		if i == 0 {
			data, _ = disc.CheckInvite(config.Friender.Invite)

			console.Log("SCC", "got guild info!", map[string]string{
				"invite": config.Friender.Invite,
			})
		}

		s, x, err := disc.JoinServer(config.Friender.Invite, data)
		if utils.HandleError(err) {
			continue
		}

		if s {
			console.Log("SCC", "joined server", map[string]string{
				"invite": config.Friender.Invite,
				"token":  token[:32] + "****",
			})
		} else {
			console.Log("DBG", "failed to join server", map[string]string{
				"reason": x,
				"token":  token[:32] + "****",
			})

			continue
		}

		if i == 0 {
			sock := discord.NewDiscordSocket(token, data.Guild.ID, data.Channel.ID)
			sock.Run()

			console.Log("SCC", "scraped", map[string]string{
				"users":     fmt.Sprint(len(sock.Users)),
				"guildID":   data.Guild.ID,
				"channelID": data.Channel.ID,
			})

			usrs = sock.Users
		}

		retries := 0
		for {
			user := usrs[xx]

			if checkIfInBlacklist(user.ID) {
				xx++

				continue
			}

			if retries == config.Friender.Retries {
				break
			}

			s, x, err := disc.SendFriendRequest(user.ID, data.Guild.ID, data.Channel.ID)
			if utils.HandleError(err) {
				continue
			}

			if s {
				sent++
				xx++

				console.Log("SCC", "sent friend request", map[string]string{
					"invite": config.Friender.Invite,
					"token":  token[:32] + "****",
					"total":  fmt.Sprint(sent),
				})

			} else {
				retries++
				console.Log("DBG", "failed to send friend request", map[string]string{
					"reason":  x,
					"token":   token[:32] + "****",
					"total":   fmt.Sprint(sent),
					"retries": fmt.Sprint(retries),
				})

				if x == "token invalid" || x == "token locked" {
					console.Log("ERR", x, map[string]string{
						"token": token[:32] + "****",
					})

					break
				}
			}

			time.Sleep(time.Second * time.Duration(config.Friender.Sleep))
		}
	}
}
