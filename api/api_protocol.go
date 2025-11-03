package api

type Command struct {
	ID         string `json:"id"`         // ResourceId + TimeStamp
	ResourceID string    `json:"resourceId"` // Cada recurso tem um ID único
	NodeID     string `json:"nodeId"`     // Ip do servidor que veio o comando
	TimeStamp  int64  `json:"timeStamp"`
	Operation  string `json:"operation"` // Create | Delete | Update
	Resource   string `json:"resource"`        // Tipo do resource que esta sendo alterado
	Value      []byte `json:"value,omitempty"` // cara resource tem um tipo diferente de dado
}

type Message struct {
	Type     string `json:"type"`              // Request, Reply, Propagate, Sync_Request,  Sync_Response
	Commands []byte `json:"command,omitempty"` // só manda command em Propagate, Sync_Request, Sync_Response |commands é um []Command
}

// Existem 3 tipos de recursos:
// Jogadores -> id Unico: Player#ID
// Partidas  -> id Unico: Match#ID
// Cartas    -> id Unico:	Card#ID
// Trocas    -> id Unico:	Trade#ID
