package main

import (
	"io/ioutil"
	"math/rand"

	"github.com/mciem/GoDMT/internal/console"
	"github.com/mciem/GoDMT/internal/discord"
	"github.com/mciem/GoDMT/internal/utils"
	"gopkg.in/yaml.v2"
)

func job(token string) {
	r := rand.Intn(len(config.Changer.Bios))
	bio := config.Changer.Bios[r]

	r1 := rand.Intn(len(config.Changer.Names))
	name := config.Changer.Names[r1]

	prox, _ := proxies.Next()
	disc, err := discord.NewDiscord(token, "http://"+prox)
	if utils.HandleError(err) {
		return
	}

	s, x, err := disc.ChangeDisplayName(name)
	if utils.HandleError(err) {
		return
	}

	if s {
		console.Log("SCC", "changed display name", map[string]string{
			"name":  name,
			"token": token[:32] + "****",
		})
	} else {
		console.Log("DBG", "failed to change display name", map[string]string{
			"reason": x,
			"token":  token[:32] + "****",
		})

		return
	}

	s1, x1, err1 := disc.ChangeBio(bio)
	if utils.HandleError(err1) {
		return
	}

	if s1 {
		console.Log("SCC", "changed bio", map[string]string{
			"bio":   name,
			"token": token[:32] + "****",
		})
	} else {
		console.Log("DBG", "failed to change bio", map[string]string{
			"reason": x1,
			"token":  token[:32] + "****",
		})

		return
	}

}

func main() {
	file, err := ioutil.ReadFile("config.yaml")
	if utils.HandleError(err) {
		return
	}

	yaml.Unmarshal(file, &config)

	tokens, _ = utils.NewFromFile("../../assets/tokens.txt")
	proxies, _ = utils.NewFromFile("../../assets/proxies.txt")

	wg.Add(len(tokens.List))

	for i := range tokens.List {
		go func(i int) {
			defer wg.Done()

			t, _ := tokens.Next()
			job(t)
		}(i)
	}

	wg.Wait()
}
