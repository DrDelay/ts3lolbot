package main

import (
	"fmt"

	"github.com/TrevorSStone/goriot"
)

type printer func(string)

func handleCommand(command string, cb printer) {
	cb(fmt.Sprintf("DEBUG: Got command %s", command))
	goriot.ChampionList(Config.Region, false)
}
