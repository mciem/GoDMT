package console

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Log(l_type string, text string, v_map map[string]string) string {
	c_time := time.Now().Format("15:04")

	if l_type == "INP" {
		fmt.Printf("%s%s %s%s %s%s > ", black, c_time, logs[l_type], l_type, white, text)

		var i_text string
		fmt.Scanln(&i_text)

		return i_text
	} else {
		var c_vars string

		i := 0
		for v, k := range v_map {
			c_vars += fmt.Sprintf(" %s%s=%s%s", black, v, white, k)

			if i != len(v_map)-1 {
				c_vars += fmt.Sprintf("%s,", black)
			}

			i++
		}

		fmt.Println(black + c_time + " " + logs[l_type] + l_type + " " + white + text + " " + c_vars)

	}

	return ""
}
