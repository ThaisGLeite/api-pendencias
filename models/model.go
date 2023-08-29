package model

type Pendencia struct {
	Dia         string  `json:"dia"`
	Nome        string  `json:"nome"`
	Description string  `json:"description"`
	Valor       float64 `json:"valor"`
}
