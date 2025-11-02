package models

type Player struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Cards    []int  `json:"cards"`
	Coins    int    `json:"coins"`
}

type Card struct {
	Id       int    `json:"id"`
	Manacost int    `json:"manacost"`
	Power    int    `json:"power"`
	Cardname string `json:"cardname"`
	Rarity   string `json:"rarity"`
}

type Trade struct {
	Id       int  `json:"id"`
	Accepted bool `json:"accepted"`
	CardA    int  `json:"cardA"`
	CardB    int  `json:"cardB"`
	PlayerA  int  `json:"playerA"`
	PlayerB  int  `json:"playerB"`
}

type Match struct {
	Id      int              `json:"id"`
	Players map[int]GameInfo `json:"players"`
}

type GameInfo struct {
	Player int    `json:"player"`
	Life   int    `json:"life"`
	Mana   int    `json:"mana"`
	Deck   [8]int `json:"deck"`
	Hand   [4]int `json:"hand"`
}
