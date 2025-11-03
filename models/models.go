package models

type Player struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Cards    []int  `json:"cards"`
	Coins    int    `json:"coins"`
}

type Card struct {
	Id       int    `json:"id"`
	Manacost int    `json:"manacost"`
	Power    int    `json:"power"`
	Quantity int    `json:"quantity,omitempty"`
	Cardname string `json:"cardname"`
	Rarity   string `json:"rarity"`
}

type Trade struct {
	Id       int    `json:"id"`
	Accepted bool   `json:"accepted"`
	PlayerA  string `json:"playerA"`
	PlayerB  string `json:"playerB,omitempty"`
	CardA    int    `json:"cardA"`
	CardB    int    `json:"cardB,omitempty"`
}

type Match struct {
	Id      int                 `json:"id"`
	Players map[string]GameInfo `json:"players"`
	Over    bool                `json:"over"`
	Turn    string              `json:"turn"` // Indica qual jogador tem a vez
}

type GameInfo struct {
	Life int    `json:"life"`
	Mana int    `json:"mana"`
	Deck [8]int `json:"deck"`
	Hand [4]int `json:"hand,omitempty"`
}
