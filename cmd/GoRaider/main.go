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

func job(token string, invite string, message string, guildID string, channelID string, massping bool) {
	pr, _ := proxies.Next()
	d, err := discord.NewDiscord(token, "http://"+pr)

	if utils.HandleError(err) {
		return
	}

	s, x, err := d.JoinServer(invite, data)
	if utils.HandleError(err) {
		return
	}

	if s {
		console.Log("SCC", "joined server", map[string]string{
			"invite": invite,
			"token":  token[:32] + "****",
		})
	} else {
		console.Log("DBG", "failed to join server", map[string]string{
			"reason": x,
			"token":  token[:32] + "****",
		})

		return
	}

	for {
		ms := message

		if massping {
			for i := 0; i < 2; i++ {
				u, _ := users.Next()
				ms += "<@" + u + ">"
			}
		}

		s, x, err := d.SendMessage(ms, guildID, channelID)
		if utils.HandleError(err) {
			continue
		}

		if s {
			console.Log("SCC", "sent message", map[string]string{
				"message": ms,
				"token":   token[:32] + "****",
			})
		} else {
			console.Log("DBG", "failed to send message", map[string]string{
				"reason": x,
				"token":  token[:32] + "****",
			})

			if x == "ratelimited" {
				time.Sleep(1 * time.Second)

				continue

			} else if x == "token locked" || x == "token invalid" {
				return
			}

			return
		}

	}
}

func scrape(token string) []string {
	ids := []string{}

	d, _ := discord.NewDiscord(token, "")
	s, x, err := d.JoinServer(config.Joiner.Invite, data)
	if utils.HandleError(err) {
		return []string{}
	}

	if s {
		console.Log("SCC", "joined server to scrape", map[string]string{
			"invite": config.Joiner.Invite,
			"token":  token[:32] + "****",
		})
	} else {
		console.Log("DBG", "failed to join server to scrape", map[string]string{
			"reason": x,
			"token":  token[:32] + "****",
		})

		return []string{}
	}
	sock := discord.NewDiscordSocket(token, data.Guild.ID, data.Channel.ID)
	sock.Run()

	console.Log("SCC", "scraped", map[string]string{
		"users":     fmt.Sprint(len(sock.Users)),
		"guildID":   data.Guild.ID,
		"channelID": data.Channel.ID,
	})

	for _, u := range sock.Users {
		ids = append(ids, u.ID)
	}

	return ids
}

func main() {
	file, err := ioutil.ReadFile("config.yaml")
	if utils.HandleError(err) {
		return
	}

	yaml.Unmarshal(file, &config)

	tokens, _ = utils.NewFromFile("../../assets/tokens.txt")
	proxies, _ = utils.NewFromFile("../../assets/proxies.txt")

	if config.Raider.Massping {
		p, _ := proxies.Next()
		t1, _ := tokens.Next()
		d, _ := discord.NewDiscord(t1, "http://"+p)
		data, _ = d.CheckInvite(config.Joiner.Invite)

		console.Log("SCC", "got guild info!", map[string]string{
			"invite": config.Joiner.Invite,
		})

		t, _ := tokens.Next()
		usrs := scrape(t)
		users = utils.New(&usrs)
	}

	wg.Add(len(tokens.List))

	for i, t := range tokens.List {
		go func(i int, t string) {
			defer wg.Done()

			job(t, config.Joiner.Invite, config.Raider.Message, config.Raider.GuildID, config.Raider.ChannelID, config.Raider.Massping)
		}(i, t)

		time.Sleep(time.Second * time.Duration(config.Joiner.Sleep))
	}

	wg.Wait()
}
