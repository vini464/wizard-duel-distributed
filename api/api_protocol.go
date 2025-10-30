package api

type Command struct {
	ID         string `json:"id"`
	NodeID     string `json:"nodeId"`
	TimeStamp  int64  `json:"timeStamp"`
	Operation  string `json:"operation"`
	Resource   string `json:"resource"`        // Tipo do resource que esta sendo alterado
	ResourceID string `json:"resourceId"`      // Cada recurso tem um ID único
	Value      []byte `json:"value,omitempty"` // cara resource tem um tipo diferente de dado
}

type Message struct {
	Type     string `json:"type"`              // Request, Reply, Propagate, Sync_Request,  Sync_Response
	Commands []byte `json:"command,omitempty"` // só manda command em Propagate, Sync_Request, Sync_Response |commands é um []Command
}

// Existem 3 tipos de recursos:
// Jogadores -> id Unico: Username
// Partidas  -> id Unico: Username1_Username2_Timestamp
// Cartas    -> id Unico: Cardname
