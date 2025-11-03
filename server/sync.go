package main

import "wizard-duel-distributed/api"

func removeDuplicates(array []api.Command) []api.Command {
	seen := make(map[string]api.Command)
	unique := []api.Command{}
	for _, e := range array {
		if c, ok := seen[e.ID]; !ok || c.TimeStamp <= e.TimeStamp {
			seen[e.ID] = e
		}
	}
	for _, c := range seen {
		unique = append(unique, c)
	}
	return unique
}

func getLatest(logs []api.Command) []api.Command {
	latests := make(map[int]api.Command) // resource: command
	for _, command := range logs {
		com, ok := latests[command.ResourceID]
		if com.Resource == command.Resource && (!ok || com.TimeStamp < command.TimeStamp) {
			latests[command.ResourceID] = command
		}
	}
	uniqueLogs := []api.Command{}
	for _, command := range logs {
		uniqueLogs = append(uniqueLogs, command)
	}
	return uniqueLogs
}
