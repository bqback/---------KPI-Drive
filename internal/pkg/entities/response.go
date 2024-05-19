package entities

type Messages struct {
	Error   *[]string `json:"error"`
	Warning *[]string `json:"warning"`
	Info    *[]string `json:"info"`
}

// Тип для анмаршалинга ответа save_Fact
type SaveResponse struct {
	Messages Messages       `json:"MESSAGES"`
	Data     map[string]int `json:"DATA"`
	Status   string         `json:"STATUS"`
}
