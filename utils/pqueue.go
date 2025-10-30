package utils

import (
	"wizard-duel-distributed/api"
)

type PriorityQueue []*api.Command

func (pq *PriorityQueue) Front() *api.Command{
	if len(*pq) > 0 {
		return (*pq)[0]
	}
	return nil
}
func (pq *PriorityQueue) Push(command api.Command) {
	p := 0
	for id, com := range *pq {
		if com.TimeStamp < command.TimeStamp {
			 p = id
		}
	}
	if p > 0 {
	prev := append((*pq)[:p], &command)
	*pq = append(prev, (*pq)[p+1:]...)
	} else {
		*pq = append(*pq, &command)
	}
}

func (pq *PriorityQueue) Pop() *api.Command {
	if len(*pq) == 0 {
		return nil
	}
	item := (*pq)[0]
	*pq = (*pq)[1:]
	return item
}
