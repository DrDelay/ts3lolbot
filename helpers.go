package main

import "strings"

func parseTsDataString(rawPayload string) *map[string]string {
	payloadPairs := strings.Split(rawPayload, " ")
	payload := make(map[string]string)
	for _, payloadPair := range payloadPairs {
		pairSplit := strings.SplitN(payloadPair, "=", 2)
		payload[pairSplit[0]] = pairSplit[1]
	}
	return &payload
}

func parseCommandString(command string) []string {
	return strings.Split(command, " ")
}
