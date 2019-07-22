package main

import (
	"fmt"
	og "stencil/ownership_generator"
)

func main() {

	users := og.GetUsersForApp("1")
	inc := 100
	for i, j := 0, inc; i < len(users) && j < len(users); i, j = j+1, j+inc {
		thread_num := j / inc
		go og.GenOwnership(users[i:j], "diaspora", "1", "mastodon", "2", thread_num)
	}
	for {
		fmt.Scanln()
	}
}
