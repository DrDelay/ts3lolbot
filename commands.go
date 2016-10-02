package main

import (
	"fmt"

	"github.com/TrevorSStone/goriot"
)

type printer func(string)

func handleCommand(command string, cb printer) {
	commandParts := parseCommandString(command)
	commandIdent := commandParts[0]

	switch commandIdent {
	case "recent":
		if len(commandParts) != 2 {
			cb("Usage: recent <summonername>")
			return
		}
		summoner := findSummoner(commandParts[1], cb)
		if summoner == nil {
			cb(fmt.Sprintf("Summoner %s not found", commandParts[1]))
			return
		}
		recent, err := goriot.RecentGameBySummoner(Config.Region, summoner.ID)
		if err != nil {
			cb("An error occured")
			fmt.Printf("API Error obtaining recent for %d: %s", summoner.ID, err.Error())
			return
		}
		for _, game := range recent {
			var result string
			if game.Statistics.Win {
				result = "Won"
			} else {
				result = "Lost"
			}
			cb(fmt.Sprintf("%s as %d - %d/%d/%d", result, game.ChampionID, game.Statistics.ChampionsKilled, game.Statistics.NumDeaths, game.Statistics.Assists))
		}
		break
	}
}

func findSummoner(name string, cb printer) *goriot.Summoner {
	summoners, err := goriot.SummonerByName(Config.Region, name)
	if err != nil {
		cb("An error occured")
		fmt.Printf("API Error searching for summ %s: %s", name, err.Error())
		return nil
	}
	for summoner := range summoners {
		summ := summoners[summoner] // Better?
		return &summ
	}
	return nil
}
