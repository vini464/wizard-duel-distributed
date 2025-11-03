package main

import (
	"encoding/json"
	"os"
	"wizard-duel-distributed/api"
)

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
	latests := make(map[string]api.Command) // resource: command
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

func UpdateLogs(filepath string, command api.Command) {
	bytes, err := os.ReadFile(filepath)
	if err == nil {
		var logs []api.Command
		err = json.Unmarshal(bytes, &logs)
		if err == nil {
			logs = append(logs, command)
			bytes, _ = json.Marshal(logs)
			file, err := os.Create(filepath)
			if err == nil {
				file.Write(bytes)
				file.Close()
				return
			}
		}
	}
}
